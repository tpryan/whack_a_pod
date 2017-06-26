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
function SCORE(){
    var total = 0;
    var pods = 0;
    var knockdowns = 0;

    this.KnockDown = function(){
        knockdowns++;
        total += 100;
    }

    this.KillPod = function(){
        pods++;
        total++;
    }

    this.GetTotal = function(){
        return total;
    }

    this.GetPods = function(){
        return pods;
    }
    this.GetKnockDowns = function(){
        return knockdowns;
    }

}

function PODS(){
    this.max = 12;
    var podsArray = [];
    var DoNotReAdd = new Object();

    var initPods = function(){
        podsArray = [];
        for (var i = 0; i < this.max; i++){
            podsArray[i] = "";
        }
    };

    this.SetMax = function(max){
        this.max = max;
        initPods();
    };
    initPods();

    this.FindEmpty = function(){
        for (var i = 0; i < this.max; i++){
            if (typeof podsArray[i] !== 'object'){
                return i;
            }
        }
        return -1;
    };

    this.IsPodPresent = function(name){
        for (var i = 0; i < podsArray.length; i++){
            if (typeof podsArray[i] !== 'object'){
                continue;
            }

            if(name == podsArray[i].name){
                return true;
            }
        }
        return false;

    };

    this.Get = function(input){
        if (typeof input == "string") {
            var name = input;
            for (var i = 0; i < podsArray.length; i++){
                if (typeof podsArray[i] !== 'object'){
                    continue;
                }

                if(name == podsArray[i].name){
                    return podsArray[i];
                }
            }
        }

        if (typeof input == "number") {
            return podsArray[input];
        }

        return ;
    };

    this.Set = function(pod){
        if (DoNotReAdd.hasOwnProperty('pod.name')){
            return;
        }
        podsArray[pod.holder] = pod;
    };

    this.Delete = function(pod){
        DoNotReAdd[pod.name] = true;
        podsArray[pod.holder] = "";
    };

    this.Count = function(){
        return podsArray.length;
    };

    this.Add = function(json){
        var name = json.metadata.name;
        var pod = new POD(json);

        if (this.IsPodPresent(name)){
            pod = this.Get(name);
        } else{
            var target = this.FindEmpty();
            pod.holder = target;
        }

        pod.SetPhase(json);

        if (pod.holder != -1){
            this.Set(pod);
        }
    }



}

function NODE(json){
    this.name = json.metadata.name;
    this.selflink = json.metadata.selfLink;
    this.type = "Node";
    this.status = "Ready";

    if (typeof json.spec.unschedulable != "undefined"){
        this.status = "Ready,SchedulingDisabled";
    }


    this.SetShortName = function(){
        var nodenameArr = this.name.split("-");
        this.shortname = nodenameArr[nodenameArr.length-1];
    }

    this.SetShortName();
}


function POD(json){
    this.name = json.metadata.name;
    this.selflink = json.metadata.selfLink;
    this.type = "Pod";
    this.host = json.spec.nodeName;
    this.hostIP = json.status.hostIP;
    this.terminateThreshold = 1000;
    this.phase = "";
    this.holder = "";
    this.shortname = "";

    if (typeof this.host == "undefined"){
        this.host = "";
    }

    this.phase ="";
    this.startTerminate ="";

    this.SetShortName = function(){
        var nodenameArr = this.name.split("-");
        this.shortname = nodenameArr[nodenameArr.length-1];
    }

    this.ShouldRemove = function(){
        if (this.phase == "terminating"){
            var now = new Date();
            if ( now - this.startTerminate > this.terminateThreshold){
                return true;
            }
        }
        return false;
    }

    this.SetPhase = function(json){
        var podPhase = json.status.phase ? json.status.phase.toLowerCase() : '';
        this.phase = podPhase;

        if ((podPhase != "terminating") && (typeof json.metadata.deletionTimestamp != "undefined")) {
            this.phase = "terminating";
            this.startTerminate = new Date();
        }
    }
    this.SetShortName();
    this.SetPhase(json);
}

function PODSUI(pods, logwindow){
    var pods = pods;
    if (typeof(logwindow)==='undefined') logwindow = new LOGWINDOW();

    var alreadyShown = new Object();
    alreadyShown.terminating = new Object();


    this.ClearTerminating = function(){
        for (var i = 0; i < pods.Count(); i++){
            var podObj = pods.Get(i);
            if (podObj.ShouldRemove() ){
                pods.Delete(podObj);
                var poddiv = document.getElementById(podObj.name);
                if (poddiv != null){
                    poddiv.parentNode.removeChild(poddiv);
                }

            }
        }
    }

    this.ClearMissing = function(podNames){
        var podsDOM = document.querySelectorAll('.pod'), i;
        for (i = 0; i < podsDOM.length; ++i) {
            if (podNames.lastIndexOf(podsDOM[i].id) < 0){
                pods.Delete(podsDOM[i].id);
            }
        }
    }

    this.ClearAll = function(){
        for (var i = 0; i < pods.Count(); i++){
            var podObj = pods.Get(i);
            var poddiv = document.getElementById(podObj.name);
            if (poddiv){
                poddiv.parentNode.removeChild(poddiv);
            }
        }
    }

    this.AddPod  = function(pod, hitHandler){


        var div = document.getElementById(pod.name);

        if (!div){
            div = document.createElement("div");
            div.id = pod.name;
            div.dataset.selflink = pod.selflink;
            div.classList.add("pod");
            var span = document.createElement("span");
            span.innerHTML = pod.shortname;
            span.dataset.selflink = pod.selflink;
            div.append(span);
            $("#pod-" + pod.holder).append(div);
            logwindow.Log(pod);
        }

        div.classList.add(pod.phase);

        if (pod.phase == "running"){
            div.addEventListener("click", hitHandler);
        } else{
            div.removeEventListener("click", hitHandler);
        }

    }

    this.DrawPods = function(json, whackHandler){

        var podNames = [];
        for (var i = 0; i < json.items.length; i++){
            podNames.push(json.items[i].metadata.name);
        }

        this.ClearTerminating();
        this.ClearMissing(podNames);

        for (var i = 0; i < json.items.length; i++){
            pods.Add(json.items[i]);
        }

        for (var i = 0; i < pods.Count(); i++){
            var pod = pods.Get(i);
            this.AddPod(pod,whackHandler);
            logwindow.Log(pod);
        }
    }
}

function API(hostname){

    this.debug = false;
    var apihostname = hostname;
    this.timeout = 2000;
    var apiprotocol = "http://"
    var uri_get = "/api/k8s/getpods?labelSelector=app%3Dapi";
    var uri_delete = "/api/k8s/deletepod/index.php?pod=";
    var uri_color = "/api/color/";
    var uri_color_complete = "/api/color-complete/";

    var ajaxProxy = function(url, successHandler, errorHandler, timeout) {
        timeout = typeof timeout !== 'undefined' ? timeout : this.timeout;
        $.ajax({
            url: url,
            success: successHandler,
            error: errorHandler,
            timeout: timeout

        });
        if (this.debug){
            console.log("Called: ", url);
        }
    };


    var getColorURI = function(){
        return apiprotocol + apihostname + uri_color;
    }

    var getColorCompleteURI = function(){
        return apiprotocol + apihostname + uri_color_complete;
    }

     var getDeleteURI = function(){
        return apiprotocol + apihostname + uri_delete;
    }

    var getPodsURI = function(){
        return apiprotocol + apihostname + uri_get;
    }

    this.Delete = function(pod, successHandler, errorHandler){
        ajaxProxy(getDeleteURI() + pod, successHandler, errorHandler);
    };

    this.Color = function(successHandler, errorHandler){
        ajaxProxy(getColorURI(), successHandler, errorHandler, 400);
    };

    this.ColorComplete = function(successHandler, errorHandler){
        ajaxProxy(getColorCompleteURI(), successHandler, errorHandler, 400);
    };

    this.Get = function(){
        ajaxProxy(getPodsURI(), successHandler, errorHandler, 500);
    };

    this.URL = getColorURI;

}

function DEPLOYMENTAPI(hostname, logwindow){
    if (typeof(logwindow)==='undefined') logwindow = new LOGWINDOW();

    this.debug = false;
    var apihostname = hostname;
    this.timeout = 2000;
    var apiprotocol = "http://"
    var uri_getnodes = "/api/k8s/getnodes";
    var uri_get = "/api/k8s/getpods?labelSelector=app%3Dapi";
    var uri_delete = "/api/k8s/deletedeploy/";
    var uri_create = "/api/k8s/createdeploy/";
    var uri_deletepod = "/api/k8s/deletepod/index.php?pod=";
    var uri_drain = "/api/k8s/drain/?node=";
    var uri_uncordon = "/api/k8s/uncordon/?node=";


    var getPodsURI = function(){
        return apiprotocol + apihostname + uri_get;
    }

    var getNodesURI = function(){
        return apiprotocol + apihostname + uri_getnodes;
    }

    var getDeleteURI = function(){
        return apiprotocol + apihostname + uri_delete;
    }

     var getDeletePodURI = function(){
        return apiprotocol + apihostname + uri_deletepod;
    }

    var getCreateURI = function(){
        return apiprotocol + apihostname + uri_create;
    }

    var getDrainURI = function(){
        return apiprotocol + apihostname + uri_drain;
    }

    var getUncordonURI = function(){
        return apiprotocol + apihostname + uri_uncordon;
    }

    var success = function(e){
        if (typeof(logwindow)!='undefined') {
            logwindow.Log(e);
        }
    }

    var error = function(e){
        console.log("Failure: " , e);
    }

    var ajaxProxy = function(url) {
        if (this.debug){
            console.log("Called: ", url);
        }
        $.ajax({
            url: url,
            success: success,
            error: error,
            timeout: this.timeout

        });

    };

    this.Delete = function(){
        ajaxProxy(getDeleteURI());
    };

    this.DeletePod = function(pod, successHandler, errorHandler){
        $.ajax({
            url: getDeletePodURI() + pod,
            success: successHandler,
            error: errorHandler,
            timeout: this.timeout

        });
    };

    this.DrainNode = function(node, successHandler, errorHandler){
        $.ajax({
            url: getDrainURI() + node,
            success: successHandler,
            error: errorHandler,
            timeout: this.timeout

        });
    };

    this.UncordonNode = function(node, successHandler, errorHandler){
        $.ajax({
            url: getUncordonURI() + node,
            success: successHandler,
            error: errorHandler,
            timeout: this.timeout

        });
    };

    this.Create = function(successHandler, errorHandler){
        $.ajax({
            url: getCreateURI(),
            success: successHandler,
            error: errorHandler,
            timeout: 1000

        });
    };

    this.Get = function(successHandler, errorHandler){
        $.ajax({
            url: getPodsURI(),
            success: successHandler,
            error: errorHandler,
            timeout: 1000

        });
        if (this.debug){
            console.log("Called: ", url);
        }
    };

    this.GetNodes = function(successHandler, errorHandler){
        $.ajax({
            url: getNodesURI(),
            success: successHandler,
            error: errorHandler,
            timeout: 1000

        });
        if (this.debug){
            console.log("Called: ", url);
        }
    };

    this.ResetNodes = function(){
        this.GetNodes(handleRefreshNodes);

    };

    var handleRefreshNodes = function(nodes){
        for (var i = 0; i < nodes.items.length; i++){
            var node = nodes.items[i];
            var url = getUncordonURI() +  node.metadata.name;
            ajaxProxy(url);
        }
    };


}

function GAME(){
    var state = "new";
    var bombShowed= false;
    var serviceDown = true;

    this.gameInterval = "";
    this.scoreInterval = "";
    this.podsInterval = "";
    this.clockInterval = "";

    this.HasBombShowed = function(){
        return bombShowed;
    }

    this.SetBombShowed = function(){
        bombShowed = true;;
    }


    this.IsServiceDown = function(){
        return serviceDown;
    }

    this.SetServiceDown = function(){
        serviceDown = true;
    }

    this.SetServiceUp = function(){
        if (this.state != "done"){
            serviceDown = false;
        }
    }


    this.GetState = function(){
        return state;
    }

    this.Start = function(colorFunction, scoreFunction, podsFunction, clockFunction){
        this.gameInterval = setInterval(colorFunction, 300);
        this.scoreInterval = setInterval(scoreFunction, 10);
        this.podsInterval = setInterval(podsFunction, 200);
        this.clockInterval = setInterval(clockFunction, 100);
        state = "started";
        startTime = Date.now();
    }

    this.Init = function(){
        console.log("Init called.")
        state = "running";
        this.SetServiceUp();
    }

    this.Stop = function(){
        state = "done";
        window.clearInterval(this.gameInterval);
        window.clearInterval(this.scoreInterval);
        window.clearInterval(this.podsInterval);
        window.clearInterval(this.clockInterval);
        this.SetServiceDown();
    }

}

function SOUNDS(){

    var hit = "";
    var hit2 = "";
    var explosion = "";
    var countdown = "";
    var startup = "";

    var makeSource = function(file,volume,loop){
        if (typeof(loop)==='undefined') loop = false;
        if (typeof(volume)==='undefined') volume = 0.3;

        var fileExt = file.split('.').pop();
        var type = "audio/wav";

        if (fileExt == "mp3"){
            type = "audio/mpeg";
        }

        var result  = new Audio();
        var src  = document.createElement("source");
        src.type = type;
        src.src  = file;
        result.preload = "auto";
        result.appendChild(src);
        result.volume = volume;
        return result;
    }

    this.SetWhack = function(filename,volume){
        hit = makeSource(filename,volume);
        hit2 = makeSource(filename,volume);
    }

    this.SetExplosion = function(filename,volume){
        explosion = makeSource(filename,volume);
    }

    this.SetCountdown = function(filename,volume){
        countdown = makeSource(filename,volume);
    }

    this.SetStartup = function(filename,volume){
        startup = makeSource(filename,volume);
    }

    this.PlayWhack = function(filename,volume){
         if (!hit.paused){
            hit2.play();
        } else{
            hit.play();
        }
    }

    this.PlayExplosion = function(filename,volume){
        explosion.play();
    }

    this.PlayCountdown = function(filename,volume){
        countdown.play();
    }

    this.PlayStartup = function(filename,volume){
        startup.play();
    }

}

function CLOCK(duration, completeHandler){
    var start_time = new Date();
    var duration = duration;
    var completeHandler = completeHandler;

    var shutItDown = function(){
        completeHandler();
        window.clearInterval(watcher);
    }

    var checkComplete = function(){
        var diff = new Date() - Date.parse(start_time);
        var count = Math.floor(diff/1000);
        if (count > duration){
            shutItDown();
        }
    }



    this.getTimeLeft = function(){
        var diff = new Date() - Date.parse(start_time);
        var count = Math.floor(diff/1000);
        var result = duration - count;
        if (result == 0){
            shutItDown();
        }
        return result;
    };
    var watcher = setInterval(checkComplete, 200);
}

function BOMBUI(waitingimg, explodeimg){
    var waiting = waitingimg;
    var explode = explodeimg;

    this.Explode = function(){
        document.querySelector("#bomb").src = explode;
        setTimeout(this.Reset, 3000);

    };

    this.Show = function(){
        $(".bombpanel").show();
    };

    this.Reset = function(){
        document.querySelector("#bomb").src = waiting;
        $(".bombpanel").hide();
        var timeToCallShow = Math.random() * 500000;
        setTimeout(this.Show, timeToCallShow);
    };

}

function LOGWINDOW(){
    var alreadyShown = new Object();
    alreadyShown.terminating = new Object();
    alreadyShown.pending = new Object();
    alreadyShown.running = new Object();


    var IsAlreadyShown = function(pod){
        if (typeof(alreadyShown[pod.phase][pod.name])==='undefined'){
            return false;
        }
        return true;
    };

    var IsError = function(e){
        if (e.kind == "Status"){
            return true;
        }
        return false;
    };

    this.Log = function(e){
        var e = jQuery.extend(true, {}, e);
        var item = e;
        if (e.kind == "Pod"){
            item = new POD(e);
        }

         if (e.kind == "Node"){
            item = new NODE(e);
        }

        if (IsError(item)){
            return;
        }

        if (item.type === "Pod"){
            if (IsAlreadyShown(item)){
                return;
            }

            alreadyShown[item.phase][item.name] = "";
            delete item.terminateThreshold;
            delete item.holder;
            if (item.startTerminate == ""){
                 delete item.startTerminate;
            }
        }


        var output = JSON.stringify(item,null,2);
        var textArray = output.split("\n");

        for (var i = textArray.length -1; i >= 0; i--){
            var css_class = "";
            var content = '<div><span>' + textArray[i] +  '</span></div>';
            if (textArray[i].indexOf("phase") >= 0){
                var css_class = "phase";
                content = '<div>  <span class="'+ css_class +'">' + textArray[i].trim() +  '</span></div>';
            }
            $(content).prependTo("#logwindow").hide().delay( (textArray.length - i) * 50 ).slideDown();
        }


    }
}