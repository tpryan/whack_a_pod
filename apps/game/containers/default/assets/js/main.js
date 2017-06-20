// Copyright 2017 Google Inc. All Rights Reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
var default_duration = 40; 

// These are set in config.js, and are specific to your cluster setup
var api = new API(servicehost);
var deploymentAPI = new DEPLOYMENTAPI(adminhost);
var pods = new PODS();
var podsUI = new PODSUI(pods);
var bombUI = new BOMBUI("assets/img/bomb_waiting.png", "assets/img/bomb_explode.png");
var game = new GAME();
var clock = "";
var score = new SCORE();
var sounds = new SOUNDS();
sounds.SetWhack("assets/audio/pop.wav",.5);
sounds.SetExplosion("assets/audio/explosion.wav",.5);
sounds.SetCountdown("assets/audio/countdown.mp3",.5);
sounds.SetStartup("assets/audio/startup.mp3",.5);

document.addEventListener('DOMContentLoaded', function() {
    $("#start-modal").show();
    $(".timer").html(default_duration);
    setReport("Kubernetes service not started yet.");
    deploymentAPI.Delete();
    var interval = Math.random() * 200000;
    document.querySelector("#bomb").addEventListener("click", bombClickHandler);
    document.querySelector("#deploy-start").addEventListener("click", startDeployment);
    document.querySelector("#restart").addEventListener("click", restart);
});

function setReport(msg, color){
    if (typeof color == "undefined") color = "#333333";
    var report = document.querySelector(".report");
    report.innerHTML = "<span>" + msg + "</span>";
    report.style.color = color;
}

function endDeployment(){
    deploymentAPI.Delete();
    game.Stop();
    showTotals();
    podsUI.ClearAll();
    setReport("Kubernetes service went away!");
    showModal("#end-modal");
}

function showTotals(){
    $("#total-pods").html(score.GetPods() + " pods");
    $("#total-knockdowns").html(score.GetKnockDowns() + " service disruptions");
    $("#total-score").html(score.GetTotal() + " points");
}

function startDeployment(){
    deploymentAPI.Create(initGame,genericError);
    hideModal("#start-modal");
    setReport("Kubernetes starting up.");
}

function initGame(){
    game.Start(getColor, showScore, getPods, getTimeLeft);
}

function getColor(){
    api.Color(handleColor, handleColorError)
}

function getTimeLeft(){
    if (clock != ""){
        $(".timer").html(clock.getTimeLeft());
        if (clock.getTimeLeft() <= 4){
            sounds.PlayCountdown();
        }
    }
}

function handleColor(e){
    if (game.GetState() == "started"){
        game.Init();
        clock = new CLOCK(default_duration, endDeployment);
        sounds.PlayStartup();
    }

    if (game.GetState() == "running"){
        game.SetServiceUp();
        setReport("Kubernetes service is UP!", e);
    }
    
}

function handleColorError(e,textStatus, errorThrown){
    if (game.GetState() == "running") {
        setReport("Kubernetes service is DOWN!", "#FF0000");
        alertYouKilledIt();
    }
}

function showScore(){
    document.querySelector(".scoreboard .total").innerHTML = score.GetTotal();
    if (score.GetTotal() >= 25 && !game.HasBombShowed()){
        bombUI.Show();
        game.SetBombShowed();
    }
}

function getPods(){
    deploymentAPI.Get(handlePods, handlePodsError);
}

function handlePods(e){
    if (game.GetState() == "done") {
        return;
    }

    podsUI.DrawPods(e, whackHandler);
}

function handlePodsError(e){
    $(".pods").html("");
    console.log("Error getting pods:", e);
}

function alertYouKilledIt(){
    if (!game.IsServiceDown() && game.GetState() == "running"){
        console.log("Killed it.");
        game.SetServiceDown();
        score.KnockDown()
        $(".alert .msg").html("You knocked down the service.");
        $(".alert").show();
        setTimeout(hideAlert, 3000);
    }
}

function whackHandler(e){
    sounds.PlayWhack();
    killPod(e.target.dataset.selflink)
}

function killPod(selflink){
    deploymentAPI.DeletePod(selflink, score.KillPod, podError);
}

function bombClickHandler(e){
    deploymentAPI.Get(bombBlastHandler, genericError);
}

function bombBlastHandler(e){
    sounds.PlayExplosion();
    for (var i = 0; i < e.items.length; i++){
        var pod = e.items[i];
        if (pod.status.phase == "Running"){
            killPod(pod.metadata.selfLink);
        }
    }
    bombUI.Explode();
}

// Add functions below to lib.js
function genericError(e){
    console.log("Error: ", e);
}

function hideModal(id){
    var modal = $(id);
    modal.fadeOut();
}

function podError(e){
    console.log("Pod already gone? :", e);
}

function showModal(id){
    var modal = $(id);
    modal.fadeIn('slow');
}

function hideAlert(){
    $(".alert").fadeOut('slow');
}