package main

import (
    "net/http"
    "strconv"
    "strings"
    "time"
    "github.com/gin-gonic/gin"
)

func CreateSupportTicket(c *gin.Context) {
    var ticket SupportTicket
    if err := c.ShouldBindJSON(&ticket); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
        return
    }
    if ticket.GoogleID == "" || ticket.Category == "" || ticket.Message == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Champs obligatoires"})
        return
    }
    ticket.Status = "ouvert"
    if err := db.Create(&ticket).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur création ticket"})
        return
    }
    c.JSON(http.StatusCreated, ticket)
}

func ListUserTickets(c *gin.Context) {
    googleID := c.Param("google_id")
    if googleID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Google ID manquant"})
        return
    }
    var tickets []SupportTicket
    err := db.Where("google_id = ?", googleID).Order("created_at desc").Find(&tickets).Error
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur BDD"})
        return
    }
    c.JSON(http.StatusOK, tickets)
}

func ListAllTickets(c *gin.Context) {
    var tickets []SupportTicket
    err := db.Order("created_at desc").Find(&tickets).Error
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur BDD"})
        return
    }
    c.JSON(http.StatusOK, tickets)
}

func ReplySupportTicket(c *gin.Context) {
    idStr := strings.TrimSpace(c.Param("id"))
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
        return
    }

    var input struct {
        AdminReply string `json:"admin_reply"`
        Status     string `json:"status"`
    }
    if err := c.ShouldBindJSON(&input); err != nil || input.AdminReply == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Réponse admin obligatoire"})
        return
    }

    var ticket SupportTicket
    if err := db.First(&ticket, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Ticket introuvable"})
        return
    }

    ticket.AdminReply = input.AdminReply
    ticket.Status = input.Status
    ticket.UpdatedAt = time.Now()

    if err := db.Save(&ticket).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de mettre à jour"})
        return
    }

    c.JSON(http.StatusOK, ticket)
}