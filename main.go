package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type UserData struct {
	Data struct {
		User struct {
			ProfilePicURLHD string `json:"profile_pic_url_hd"`
		} `json:"user"`
	} `json:"data"`
}

func downloadProfilePic(username string) error {
	url := fmt.Sprintf("https://i.instagram.com/api/v1/users/web_profile_info/?username=%s", username)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 11)")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-IG-App-ID", "936619743392459") 

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Instagram returned %d: %s", resp.StatusCode, string(body))
	}

	var data UserData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("JSON decode failed: %v", err)
	}

	picURL := data.Data.User.ProfilePicURLHD
	if picURL == "" {
		return fmt.Errorf("could not find profile picture (maybe private account?)")
	}
	resp, err = http.Get(picURL)
	if err != nil {
		return fmt.Errorf("failed to download picture: %v", err)
	}
	defer resp.Body.Close()

	file, err := os.Create(username + ".jpg")
	if err != nil {
		return fmt.Errorf("file create failed: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func main() {
	username := "Hacker_jai_op_pvt" // username
	if err := downloadProfilePic(username); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Profile picture downloaded successfully!")
	}
}