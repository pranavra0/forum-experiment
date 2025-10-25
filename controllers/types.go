package controllers

import "forum-experiment/models"

type PageData struct {
	Name     string
	Threads  []models.Thread
	Replies  []models.Reply
}