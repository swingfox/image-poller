<h3 align="center">Image Poller</h3>

---

## 📝 Table of Contents

- [About](#about)
- [Getting Started](#getting_started)
- [Deployment](#deployment)
- [API Usage](#usage)
- [Built Using](#built_using)

## 🧐 About <a name = "about"></a>

The goal of this project is to create a consumable RESTful API for storing and retrieving images. The application will call Pexel to retrieve random images and save it to the Cloudinary for image data store.
The application uses API keys to communicate with Pexel and Cloudinary. 

## 🏁 Getting Started <a name = "getting_started"></a>

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See [deployment](#deployment) for notes on how to deploy the project on a live system.

### Prerequisites

What things you need to install the software and how to install them.

```
Go Lang 1.18
Go Land
MongoDB 5.0.9
```

## 🎈 API Usage <a name="usage"></a>

## GET /images

Get images from the Pexel and saves it to the Cloudinary.

## GET /images/{id}

Get image metadata by ID.

## PATCH /images/{id}

Update certain fields in the images metadata.

## DELETE /images/{id}

Perform a soft delete of the image.

## 🚀 Deployment <a name = "deployment"></a>

Execute command below to run the whole application.

```
go run .\cmd\webapp\main.go
```

## ⛏️ Built Using <a name = "built_using"></a>

- [MongoDB](https://www.mongodb.com/) - Database
- [Go Lang 1.18](https://go.dev/blog/go1.18)
- [GoLand](https://www.jetbrains.com/go/)
