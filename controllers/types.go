package controllers

import "forum-experiment/models"

type PageData struct {
	Name     string
	Threads  []models.Thread
	Replies  []models.Reply
	User *models.User
	Error    string
}

