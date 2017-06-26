<?php
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
    header("Access-Control-Allow-Origin: *");

    function getToken(){
        return file_get_contents("/var/run/secrets/kubernetes.io/serviceaccount/token");
    }

    function getK8sCurlHandle(){
        $kube_token = getToken();
        $ch = curl_init();
        $headers = ["Authorization: Bearer " . $kube_token];


        curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
        curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);

        return $ch;
    }

    function getK8sCurlHandleForPost($input){
        $kube_token = getToken();
        $ch = curl_init();

        curl_setopt($ch, CURLOPT_HTTPHEADER, array(
            'Content-Type: application/json',
            'Content-Length: ' . strlen($input),
            'Authorization: Bearer ' . $kube_token)
        );
        curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);

        return $ch;
    }

    function getK8sCurlHandleForPatch($input){
        $kube_token = getToken();
        $ch = curl_init();

        curl_setopt($ch, CURLOPT_HTTPHEADER, array(
            'Content-Type: application/merge-patch+json',
            'Content-Length: ' . strlen($input),
            'Authorization: Bearer ' . $kube_token)
        );
        curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);

        return $ch;
    }

    function killAllPods(){
        $ch = getK8sCurlHandle();
        curl_setopt($ch, CURLOPT_CUSTOMREQUEST, "GET");
        curl_setopt($ch, CURLOPT_URL, "https://kubernetes/api/v1/pods?labelSelector=app=api");
        $output = curl_exec($ch);
        $podList = json_decode($output);
        $pods = $podList->items;

        $podsToPurge = array();
        for ($i = 0; $i < count($pods); $i++){
            array_push($podsToPurge, $pods[$i]->metadata->selfLink);
        }

        curl_setopt($ch, CURLOPT_CUSTOMREQUEST, "DELETE");
        for ($i = 0; $i < count($podsToPurge); $i++){
            curl_setopt($ch, CURLOPT_URL, "https://kubernetes" . $podsToPurge[$i]);
            $output = curl_exec($ch);
        }
        curl_close($ch);

    }

    function killAllPodsOnNode($node){
        $ch = getK8sCurlHandle();
        curl_setopt($ch, CURLOPT_CUSTOMREQUEST, "GET");
        curl_setopt($ch, CURLOPT_URL, "https://kubernetes/api/v1/pods?labelSelector=app=api");
        $output = curl_exec($ch);
        $podList = json_decode($output);
        $pods = $podList->items;

        $podsToPurge = array();
        for ($i = 0; $i < count($pods); $i++){
            if ($pods[$i]->spec->nodeName == $node){
                array_push($podsToPurge, $pods[$i]->metadata->selfLink);
            }
        }

        curl_setopt($ch, CURLOPT_CUSTOMREQUEST, "DELETE");
        for ($i = 0; $i < count($podsToPurge); $i++){
            curl_setopt($ch, CURLOPT_URL, "https://kubernetes" . $podsToPurge[$i]);
            curl_exec($ch);
        }
        curl_close($ch);
    }

