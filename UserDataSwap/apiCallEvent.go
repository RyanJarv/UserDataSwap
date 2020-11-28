package main

type InstanceState struct {
	Name string `json:"name"`
}

type UserData struct {
	Value string `json:"Value"`
}

type InstanceItems struct {
	InstanceId    string        `json:"instanceId,omitempty"`
	InstanceState InstanceState `json:"instanceState,omitempty"`
}

type InstanceSet struct {
	Items []InstanceItems `json:"items"`
}

type ResponseElements struct {
	InstancesSet InstanceSet `json:"instancesSet"`
}

type RequestParameters struct {
	InstancesSet InstanceSet `json:"instancesSet"`
	UserData string `json:"userData,omitempty"`
}

type RunInstancesEvent struct {
	EventVersion string `json:"eventVersion"`
	EventName string `json:"eventName"`
	AwsRegion string `json:"awsRegion"`
	RequestParameters RequestParameters `json:"requestParameters"`
	ResponseElements ResponseElements `json:"responseElements"`
}
