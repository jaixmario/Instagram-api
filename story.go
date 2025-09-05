package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func getUserID(username, sessionID string) (string, error) {
	url := fmt.Sprintf("https://i.instagram.com/api/v1/users/web_profile_info/?username=%s", username)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Instagram 155.0.0.37.107")
	req.Header.Set("Cookie", "sessionid="+sessionID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch user id, status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	user := result["data"].(map[string]interface{})["user"].(map[string]interface{})
	return user["id"].(string), nil
}

func downloadStories(userID, sessionID string) error {
	url := fmt.Sprintf("https://i.instagram.com/api/v1/feed/user/%s/story/", userID)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Instagram 155.0.0.37.107")
	req.Header.Set("Cookie", "sessionid="+sessionID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to fetch stories, status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	items := result["reel"].(map[string]interface{})["items"].([]interface{})
	for i, item := range items {
		story := item.(map[string]interface{})
		var url string
		if video, ok := story["video_versions"]; ok {
			url = video.([]interface{})[0].(map[string]interface{})["url"].(string)
		} else {
			url = story["image_versions2"].(map[string]interface{})["candidates"].([]interface{})[0].(map[string]interface{})["url"].(string)
		}

		fmt.Println("Downloading:", url)
		if err := saveFile(fmt.Sprintf("story_%d.mp4", i), url, sessionID); err != nil {
			return err
		}
	}
	return nil
}

func saveFile(filename, url, sessionID string) error {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", "sessionid="+sessionID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, _ := os.Create(filename)
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func main() {
	username := "the.rebel.kid"
	sessionID := "SESSION ID FROM COOKIES"

	userID, err := getUserID(username, sessionID)
	if err != nil {
		fmt.Println("Error fetching user ID:", err)
		return
	}
	fmt.Println("User ID:", userID)

	if err := downloadStories(userID, sessionID); err != nil {
		fmt.Println("Error downloading stories:", err)
	}
}