package command

import (
	"encoding/json"
	"fmt"
	"strings"
)

// vehicleCategory turns a vehicle name into the index the API keys its pricing
// model by. Getting it wrong prices a shipment against another vehicle, so it
// errors instead of silently falling back to bici.
func vehicleCategory(name string) (int, error) {
	switch name {
	case "bici", "moto":
		return int(VehicleCategoryBici), nil
	case "car", "auto":
		return int(VehicleCategoryCar), nil
	case "van", "miniflete":
		return int(VehicleCategoryVan), nil
	case "truck", "camion":
		return int(VehicleCategoryTruck), nil
	case "rover":
		return int(VehicleCategoryRover), nil
	case "drone":
		return int(VehicleCategoryDrone), nil
	}
	return 0, fmt.Errorf("unknown vehicle: %s (bici, car, van, truck, rover, drone)", name)
}

// unwrap pulls a single value out of a response envelope: the API answers with
// {"request_id": "..."} and friends.
func unwrap(body, key string) (string, error) {
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		return "", err
	}

	value, ok := response[key].(string)
	if !ok {
		return "", fmt.Errorf("no %s in response: %s", key, body)
	}

	return value, nil
}

// geocodeAddress turns a free-form address into the Address shape the API stores,
// reusing the platform's own geocoder so the coordinates match what it routes on.
func geocodeAddress(address string) (map[string]interface{}, error) {
	formatted, latitude, longitude, err := latlng(address)
	if err != nil {
		return nil, err
	}

	countryCode := "UY"
	if strings.Contains(formatted, "Argentina") {
		countryCode = "AR"
	}

	return map[string]interface{}{
		"address_lines": map[string]interface{}{"0": formatted},
		"country_code":  countryCode,
		"latitude":      latitude,
		"longitude":     longitude,
	}, nil
}

// whoAmI returns the caller's user ID.
func whoAmI() (string, error) {
	body, err := GET(fmt.Sprintf("%s/me?authorization=%s", getUri(), getToken()), "me")
	if err != nil {
		return "", err
	}

	var response map[string]interface{}
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		return "", err
	}

	user, ok := response["user"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("no user in response: %s", body)
	}

	id, ok := user["id"].(string)
	if !ok {
		return "", fmt.Errorf("no id in user: %s", body)
	}

	return id, nil
}

func findOperation(operationID string) (map[string]interface{}, error) {
	body, err := GET(fmt.Sprintf("%s/operations/%s?authorization=%s", getUri(), operationID, getToken()), "operation")
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		return nil, err
	}

	operation, ok := response["operation"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no operation in response: %s", body)
	}

	return operation, nil
}
