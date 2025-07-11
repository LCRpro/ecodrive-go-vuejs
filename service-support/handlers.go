package main

import (
    "net/http"
    "strconv"
    "strings"
    "time"
    "github.com/gin-gonic/gin"
)

var tickets []SupportTicket
var nextTicketID uint = 1

func saveTicket(ticket *SupportTicket) {
    if ticket.ID == 0 {
        ticket.ID = nextTicketID
        nextTicketID++
    }
    if ticket.CreatedAt.IsZero() {
        ticket.CreatedAt = time.Now()
    }
    tickets = append(tickets, *ticket)
}

func findTicketsByGoogleID(googleID string) []SupportTicket {
    var res []SupportTicket
    for _, t := range tickets {
        if t.GoogleID == googleID {
            res = append(res, t)
        }
    }
    return res
}

func findTicketByID(id uint) (*SupportTicket, bool) {
    for i := range tickets {
        if tickets[i].ID == id {
            return &tickets[i], true
        }
    }
    return nil, false
}

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
    saveTicket(&ticket)
    c.JSON(http.StatusCreated, ticket)
}

func ListUserTickets(c *gin.Context) {
    googleID := c.Param("google_id")
    if googleID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Google ID manquant"})
        return
    }
    ticketsFound := findTicketsByGoogleID(googleID)
    c.JSON(http.StatusOK, ticketsFound)
}

func ListAllTickets(c *gin.Context) {
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

    ticket, ok := findTicketByID(uint(id))
    if !ok {
        c.JSON(http.StatusNotFound, gin.H{"error": "Ticket introuvable"})
        return
    }

    ticket.AdminReply = input.AdminReply
    ticket.Status = input.Status
    ticket.UpdatedAt = time.Now()


    c.JSON(http.StatusOK, ticket)
}