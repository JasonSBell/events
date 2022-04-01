# Events

This repository contains the code for Allokate's Events service. This service provides REST API endpoints for publishing and managing events that happen within the application.

# Table of Content

- [Events](#events)
- [Table of Content](#table-of-content)
- [Overview](#overview)
  - [Event structure](#event-structure)
- [Endpoints](#endpoints)

# Overview

The role of this service is to simplify the operations of the various other components and provide a unified means of publishing events so that an event sourcing strategy may be employed within the app to orchestrate and chain together various tasks and jobs that should run only when certain events happen. For example:

- An record of an account's/portfolio's holdings may need to be updated if new transactions are imported
- The ticker scanner may need to be run against a newly created article to find stock ticker symbols.
- The sentiment analysis process may need to be run against a newly created article to classify the sentiment toward the various stocks within the article.

The service writes semi structured event data to the Mongo database and subsequently routes events as specified by the provided configuration file. The configuration file and the service currently allows events to be routed to various exchanges and queues within the RabbitMQ service based on the event type. The service and associated configuration file also allows events to be routed to a user via the notification service.

## Event structure

```json
{
  "timestamp": "2022-03-19T12:31:52Z",
  "type": "article.created",
  "body": {
    // Specific to the event this field can be anything from a number, to a string, to a JSON body or a binary blob.
    "id": "622df6275e872cc36c45056a",
    "title": " New Tax Rules Force Faster Payouts for Some IRA Holders"
  },
  "source": "articles" // The name of the service or origin of the event.
}
```

# Endpoints

- GET /api/ping
- PUT /api/events
- GET /api/events
- GET /api/events/:id
- DELETE /api/events/:id
- GET /api/events/types
