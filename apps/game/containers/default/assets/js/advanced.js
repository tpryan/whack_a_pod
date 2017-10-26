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

var api = new API(servicehost);
var logwindow = new LOGWINDOW();
var deploymentAPI = new DEPLOYMENTAPI(adminhost, logwindow);
var pods = new PODS();
var game = new GAME();
var score = new SCORE();
var nodes = [];
var pods_shown = [];
var fails_threshold = 20;


document.addEventListener('DOMContentLoaded', function() {
    api.timeout = 10000;
    deploymentAPI.Delete();
    deploymentAPI.ResetNodes();
    setReport("");
    $("#deploy-start").click(startDeployment);
    $("#deploy-end").click(endDeployment);
    $("#endpoint").html(api.URL());
    $("#show-pod-yaml").click(showPodModal);
    $("#close-pod-modal").click(hidePodModal);
    $("#show-service-yaml").click(showServiceModal);
    $("#close-service-modal").click(hideServiceModal);
});

function showModal(id){
    var modal = $(id);
    modal.fadeIn('slow');
}

function hideModal(id){
    var modal = $(id);
    modal.fadeOut();
}

function showPodModal(e){
    showModal("#pod-yaml");
}

function hidePodModal(e){
    hideModal("#pod-yaml");
}

function showServiceModal(e){
    showModal("#service-yaml");
}

function hideServiceModal(e){
    hideModal("#service-yaml");
}

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

function startDeployment(){
    $("#deploy-start").hide();
    $("#deploy-start").css('z-index', 1);
    $("#deploy-end").show();
    $("#deploy-end").css('z-index', 300);
    deploymentAPI.Create(initGame,genericError);
    setReport("");
    $("#deployment").html("kubectl delete -f whack-a-pod-deployment.yaml")
}

function initGame(e){
    game.Start(getColor, showScore, getPods, getTimeLeft);
    logwindow.Log(e);

}

function getColor(){
    api.ColorComplete(handleColor, handleColorError)
}

function handleColor(e){
    if (game.GetState() == "started"){
        game.Init();
        setTimeout(function(){$(".pods-explain").fadeOut();}, 3000);
    }

    if (game.GetState() == "running"){
        game.SetServiceUp();
        setReport("Service call result: "+ e.color, e.color);
        $(".responder").removeClass("responder");
        $("#"+e.name).addClass("responder");

    }
    if (api.ResetFails()){
        console.log("Soft service has recovered.");
    }
}

function handleColorError(e,textStatus, errorThrown){
    if (game.GetState() == "running") {
        if (api.IsHardFail()){
            console.log("Hard service fail.");
            setReport("Kubernetes service is DOWN!", "#FF0000");
            $(".responder").removeClass("responder");
            alertYouKilledIt();
        } else {
            console.log("Soft service fail. Retry");
        }
    }
}

function getPods(){
    deploymentAPI.Get(handlePods, handlePodsError);
}

function handlePods(e,b){
    
    if (game.GetState() == "done") {
        return;
    }

    drawPods(e);
    
}

function drawPods(pods){

    var pods_active = [];
    //create node UI if it doesn't exist
    for (var i = 0; i < pods.items.length; i++){
        var pod = new POD(pods.items[i]);
        logwindow.Log((pod));

        if ((pod.host == null) || (pod.host.length == 0)){
            continue;
        }

        pods_active.push(pod.name);

        if (nodes.indexOf(pod.host) < 0){
            nodes.push(pod.host);
            createNodeUI(pod.host);
        }

        if (pods_shown.indexOf(pod.name) < 0){
            pods_shown.push(pod.name);
            createPodUI(pod);
        }

        var $pod = $("#"+ pod.name);
        $pod.addClass(pod.phase);

        if (pod.phase == "running") {
            if (!$pod.hasClass("bound")){
                $pod.click(whackHandler);
                $pod.addClass("bound");
            }
        } else{
            $pod.click();
        }

        if (pod.ShouldRemove()){
            $pod.remove();
        }

    }
    for (var i = 0; i < pods_shown.length; i++){
        if (pods_active.indexOf(pods_shown[i]) < 0){
            $("#"+ pods_shown[i]).remove();
        }
    }

}


function createPodUI(pod){
    var hostID = "node_"+pod.host;

    var div = document.createElement("div");
    div.id = pod.name;
    div.dataset.selflink = pod.selflink;
    div.classList.add("pod");

    var span = document.createElement("span");
    span.innerHTML = pod.shortname;
    span.dataset.selflink = pod.selflink;
    div.append(span);

    var label = document.createElement("div");
    label.classList.add("kube-label");
    label.classList.add("kube-pod");
    label.innerHTML= "Pod";
    div.append(label);


    $("#pod-" + pod.holder).append(div);
    $("#"+ hostID).append(div);

    // logwindow.Log((pod));

}


function createNodeUI(name){
    var id = "node_"+name;
    var label = name.split("-")[name.split("-").length -1 ];
    var $div = $('<div class="node" id="' + id + '"><button class="small" id="kill-' + id +  '">X</button><p>Node: <strong>' + label + '</strong></p><div class="kube-label kube-node">Node</div></div>');
    var $holder = $("#nodes");
    $div.appendTo("#nodes")
    $killbtn = $("#kill-" + id)
    
    if(name=="minikube"){
        $killbtn.hide();
    } else{
        $killbtn.click(killNode);
        $killbtn.data("node", name);
    }
    
}

function killNode(e){
    var $killbtn = $("#" + e.currentTarget.id);
    var node = $killbtn.data("node");
    deploymentAPI.DrainNode(node);
    $killbtn.click(resetNode);
    $killbtn.addClass("reset");
    $killbtn.text("+");
}


function resetNode(e){
    var $killbtn = $("#" + e.currentTarget.id);
    var node = $killbtn.data("node");
    deploymentAPI.UncordonNode(node);
    $killbtn.removeClass("reset");
    $killbtn.text("X");
    $killbtn.click(killNode);

}

function handlePodsError(e){
    // $(".pods").html("");
    if (typeof e != "string"){
        console.log("Error getting pods:", e.statusText);
    } else{
        console.log("Error getting pods:", e);
    }
    
}

function genericError(e){
    console.log("Error: ", e);
}

function showScore(){
}

function getTimeLeft(){
}

function alertYouKilledIt(){
    if (!game.IsServiceDown() && game.GetState() == "running"){
        console.log("Killed it.");
        game.SetServiceDown();
        score.KnockDown()
    }
}


function whackHandler(e){
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
    logwindow.Log(e);
}

function podError(e){

    // console.log("Pod already gone? :", e);
}

function endDeployment(){
    deploymentAPI.Delete();
    game.Stop();
    setReport("");
    location.reload();
}
