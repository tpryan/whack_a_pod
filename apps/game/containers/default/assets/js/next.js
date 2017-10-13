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
var logwindow = new LOGWINDOW();
var deploymentAPI = new DEPLOYMENTAPI(adminhost, logwindow);
var pods = new PODS();
var podsUI = new PODSUI(pods, logwindow);
var bombUI = new BOMBUI("assets/img/bomb_waiting_next.png", "assets/img/bomb_explode_next.png");
var game = new GAME();
var clock = "";
var score = new SCORE();
var sounds = new SOUNDS();
var fails_threshold = 9;
sounds.SetWhack("assets/audio/pop.wav",.5);
sounds.SetExplosion("assets/audio/explosion.wav",.5);
sounds.SetCountdown("assets/audio/countdown.mp3",.5);
sounds.SetStartup("assets/audio/startup.mp3",.5);


document.addEventListener('DOMContentLoaded', function() {
    $("#start-modal").show();
    $(".timer").html(default_duration);
    setReport("");
    deploymentAPI.Delete();
    var interval = Math.random() * 200000;
    $("#bomb").click(bombClickHandler);
    $("#deploy-start").click(startDeployment);
    $("#restart").click(restart);
    $("#deploy-start").focus();
});


function setReport(msg, color){
    if (typeof color == "undefined") color = "#333333";
    var report = $(".report");
    
    report.css("-webkit-filter", "drop-shadow(2px 2px 3px " + color + ")");
    report.css("color", color);
    var msgholder = $("#report_message");
    if (msgholder.length == 0){
        report.after('<div id="report_message">' + msg + '</div>')
    } else {
        msgholder.html(msg);
    }


    if(game.GetState() == "running"){
        if(game.IsServiceDown()) {
            $(".report").addClass("service_down");
            $(".report").removeClass("service_up");
        } else{
            $(".report").addClass("service_up");
            $(".report").removeClass("service_down");
        }
    } 
    
}

function endDeployment(){
    deploymentAPI.Delete();
    game.Stop();
    showTotals();
    podsUI.ClearAll();
    setReport("");
    $(".report").removeClass("service_up");
    $(".report").removeClass("service_down");
    showModal("#end-modal");
    setTimeout(function(){location.reload();}, 15000);
}

function showTotals(){
    $("#total-pods").html(score.GetPods() + " pods");
    $("#total-knockdowns").html(score.GetKnockDowns() + " service disruptions");
    $("#total-score").html(score.GetTotal() + " points");
}

function startDeployment(){
    deploymentAPI.Create(initGame,genericError);
    hideModal("#start-modal");
    $(".pods-explain").fadeIn();
    
    setReport("");
}

function initGame(e){
    game.Start(getColor, showScore, getPods, getTimeLeft);
    logwindow.Log(e);
    
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
        setTimeout(function(){$(".pods-explain").fadeOut();}, 3000);
    }

    if (game.GetState() == "running"){
        game.SetServiceUp();
        setReport("", e);
    }
    if (api.fails > 0){
        api.fails = 0;
        console.log("Soft service has recovered.");
    }
    
    
}

function handleColorError(e,textStatus, errorThrown){
    if (game.GetState() == "running") {
        if (api.fails > fails_threshold){
            console.log("Hard service fail.");
            setReport("Kubernetes service is DOWN!", "#FF0000");
            alertYouKilledIt();
        } else {
            console.log("Soft service fail. Retry");
            api.fails++;
        }
        
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
    if (e.target.id == ""){
        $("#" + e.target.parentNode.id ).addClass("terminating");
    } else{
        $("#" + e.target.id ).addClass("terminating");
    }
    
   
    killPod(e.target.dataset.selflink)
}

function killPod(selflink){
    deploymentAPI.DeletePod(selflink, killHandler, podError);
}

function killHandler(e){
    score.KillPod();
    logwindow.Log(e);
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


function showModal(id){
    var modal = $(id);
    modal.fadeIn('slow');
}

function hideModal(id){
    var modal = $(id);
    modal.fadeOut();
}

function restart(){
    location.reload();
}

function hideAlert(){
    $(".alert").fadeOut('slow');
}

function podError(e){
    console.log("Pod already gone? :", e);
}

function genericError(e){
    console.log("Error: ", e);
}