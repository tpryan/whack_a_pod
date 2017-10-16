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
	include "../lib.php";


	$body = file_get_contents(dirname(__DIR__) . "/createdeploy/deployment.json");

	if(isset($_GET['replicas'])) {
		$body = str_replace ( '"replicas": 12,' , '"replicas": '. $_GET['replicas'] .',' ,$body);
	}


	$ch = getK8sCurlHandleForPost($body);
	curl_setopt($ch, CURLOPT_CUSTOMREQUEST, "POST");
	curl_setopt($ch, CURLOPT_URL, "https://kubernetes/apis/extensions/v1beta1/namespaces/default/deployments");
	curl_setopt($ch, CURLOPT_POSTFIELDS, $body);

	$output = curl_exec($ch);
	$json = json_decode($output);

	$deploymentResult = array();
	$deploymentResult['item'] = "DEPLOYMENT";
	$deploymentResult['selflink'] = $json->metadata->selfLink;
	if (isset($json->code) && $json->code == "409"){
		$deploymentResult['status'] = "ALREADY EXISTS";
	} else{
		$deploymentResult['status'] = "CREATED";
	}

	curl_close($ch);

	header("Content-Type: application/json;charset=utf-8");
	echo json_encode($deploymentResult);