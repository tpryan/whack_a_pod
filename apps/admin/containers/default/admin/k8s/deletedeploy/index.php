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
	$ch = getK8sCurlHandle();
	$result = array();

	$deploymentSelfLink = "/apis/extensions/v1beta1/namespaces/default/deployments/api-deployment";

	curl_setopt($ch, CURLOPT_URL, "https://kubernetes/apis/extensions/v1beta1/replicasets");
	$output = curl_exec($ch);
	$replicaSets = json_decode($output);
	$sets = $replicaSets->items;

	$replicaSetSelfLink = "None found";
	for ($i = 0; $i < count($sets); $i++){
		$set = $sets[$i];
		if (isset($set->metadata->labels->app) && $set->metadata->labels->app == "api"){
			$replicaSetSelfLink = $set->metadata->selfLink;
			break;
		}
	}


	// Kill deployment
	curl_setopt($ch, CURLOPT_CUSTOMREQUEST, "DELETE");
	curl_setopt($ch, CURLOPT_URL, "https://kubernetes" . $deploymentSelfLink);
	$output = curl_exec($ch);
	$json = json_decode($output);


	$deploymentResult = array();
	$deploymentResult['item'] = "DEPLOYMENT";
	$deploymentResult['selflink'] = $deploymentSelfLink;

	if (isset($json->code) &&  $json->code == "404"){
		$deploymentResult['status'] = "NOT FOUND";
	} else{
		$deploymentResult['status'] = "DELETED";
	}


	array_push($result, $deploymentResult);


	//Kill replica set
	curl_setopt($ch, CURLOPT_CUSTOMREQUEST, "DELETE");
	curl_setopt($ch, CURLOPT_URL, "https://kubernetes" . $replicaSetSelfLink);
	$output = curl_exec($ch);
	$json = json_decode($output);

	$replicaSetResult = array();
	$replicaSetResult['item'] = "REPLICASET";
	$replicaSetResult['selflink'] = $replicaSetSelfLink;

	if (isset($json->code) &&  $json->code == "404"){
		$replicaSetResult['status'] = "NOT FOUND";
	} else{
		$replicaSetResult['status'] = "DELETED";
	}

	array_push($result, $replicaSetResult);

	curl_setopt($ch, CURLOPT_CUSTOMREQUEST, "GET");
	curl_setopt($ch, CURLOPT_URL, "https://kubernetes/api/v1/pods?labelSelector=app=api");
	$output = curl_exec($ch);
	$podList = json_decode($output);
	$pods = $podList->items;


	$podsToPurge = array();
	for ($i = 0; $i < count($pods); $i++){
		array_push($podsToPurge, $pods[$i]->metadata->selfLink);
	}

	for ($i = 0; $i < count($podsToPurge); $i++){
		curl_setopt($ch, CURLOPT_CUSTOMREQUEST, "DELETE");
		curl_setopt($ch, CURLOPT_URL, "https://kubernetes" . $podsToPurge[$i]);
		$output = curl_exec($ch);

		$podResult = array();
		$podResult['item'] = "POD";
		$podResult['selflink'] = $podsToPurge[$i];
		$podResult['status'] = "DELETED";
		array_push($result, $podResult);

	}
	curl_close($ch);

	header("Content-Type: application/json;charset=utf-8");
	echo json_encode($result);

?>
