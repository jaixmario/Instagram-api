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
			Username        string `json:"username"`
			FullName        string `json:"full_name"`
			Biography       string `json:"biography"`
			IsPrivate       bool   `json:"is_private"`
			ProfilePicURLHD string `json:"profile_pic_url_hd"`
		} `json:"user"`
	} `json:"data"`
}

func fetchUserDetails(username string) (*UserData, error) {
	url := fmt.Sprintf("https://i.instagram.com/api/v1/users/web_profile_info/?username=%s", username)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 11)")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-IG-App-ID", "936619743392459") // public IG app ID

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Instagram returned %d: %s", resp.StatusCode, string(body))
	}

	var data UserData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("JSON decode failed: %v", err)
	}

	return &data, nil
}

func downloadProfilePic(url, username string) error {
	resp, err := http.Get(url)
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
	username := "Hacker_jai_op" //username

	userData, err := fetchUserDetails(username)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	user := userData.Data.User

	fmt.Println("📌 Instagram Profile Info")
	fmt.Println("Username:", user.Username)
	fmt.Println("Name:", user.FullName)
	fmt.Println("Bio:", user.Biography)
	if user.IsPrivate {
		fmt.Println("Account Type: 🔒 Private")
	} else {
		fmt.Println("Account Type: 🌍 Public")
	}
	fmt.Println("Profile Pic URL:", user.ProfilePicURLHD)

	// Download profile pic
	if err := downloadProfilePic(user.ProfilePicURLHD, username); err != nil {
		fmt.Println("Error downloading picture:", err)
	} else {
		fmt.Println("✅ Profile picture saved as", username+".jpg")
	}
}