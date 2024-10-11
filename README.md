![enter image description here](https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Flogos-world.net%2Fwp-content%2Fuploads%2F2020%2F11%2FTripadvisor-Logo.png&f=1&nofb=1&ipt=33e5e7ec0e5cb0e4fe17cfd6558703db48ccbe6bbaf337b69491a90e6a82f69e&ipo=images)

# Tripadvisor Reviews Scrapper

## About This Actor
This Actor is a powerful, user-fiendly tool made to scrape reviews from specified Tripadvisor listings. This tool will save you time and provide you with reliable data on reviews about your provided listings.

Made with Golang 1.22.1

## Tutorial

Basic Usage

```json
{
    "startUrls": [
        {
            "url": "https://www.tripadvisor.com/Hotel_Review-g188107-d231860-Reviews-Beau_Rivage_Palace-Lausanne_Canton_of_Vaud.html"
        }
    ],
    "offset": 0,
    "limit": 20
}
```

| parameter | type | argument | description |

| --------- | ----- | ------------------------- | ---------------------------- |

| startUrls | array | _[]_ | An array of start urls |

| offset | int | _default=0_ | Start from a specific offset |

| limit | int | _default=20_ | Limit number of results |

### Output Sample

```json

[
     {
    "id": 972086842,
    "createdDate": "2024-09-30",
    "publishedDate": "2024-09-30",
    "rating": 5,
    "publishPlatform": "OTHER",
    "tripInfo": {
      "stayDate": "2023-11-30",
      "tripType": "BUSINESS"
    },
    "photoIds": [],
    "locationId": 231860,
    "labels": [],
    "title": "Super Awesome Review",
    "text": "This place was amazing! Would come again",
    "url": "/ShowUserReviews-g188107-d231860-r972086842-Beau_Rivage_Palace-Lausanne_Canton_of_Vaud.html",
    "photos": [],
    "userProfile": {
      "isMe": false,
      "isVerified": false,
      "contributionCounts": {
        "sumAllUgc": 3,
        "sumAllLikes": 0
      },
      "isFollowing": false,
      "id": "F76EE9C50E5455BA668AEE543412352A",
      "userId": "F76EE9C50E5455BA668AEE543412352A",
      "displayName": "Super Awesome",
      "username": "superReviewer",
      "hometown": {
        "locationId": null,
        "location": null,
        "fallbackString": null
      },
      "route": {
        "url": "/Profile/superReviewer"
      },
      "avatar": {
        "id": 452390350,
        "photoSizes": [
          {
            "width": 0,
            "height": 0,
            "url": "https://dynamic-media-cdn.tripadvisor.com/media/photo-o/1a/f6/ed/ce/default-avatar-2020-7.jpg?w=100&h=100&s=1"
          },
          {
            "width": 50,
            "height": 50,
            "url": "https://media-cdn.tripadvisor.com/media/photo-t/1a/f6/ed/ce/default-avatar-2020-7.jpg"
          },
          {
            "width": 150,
            "height": 150,
            "url": "https://media-cdn.tripadvisor.com/media/photo-l/1a/f6/ed/ce/default-avatar-2020-7.jpg"
          },
          {
            "width": 180,
            "height": 200,
            "url": "https://media-cdn.tripadvisor.com/media/photo-i/1a/f6/ed/ce/default-avatar-2020-7.jpg"
          },
          {
            "width": 205,
            "height": 205,
            "url": "https://media-cdn.tripadvisor.com/media/photo-f/1a/f6/ed/ce/default-avatar-2020-7.jpg"
          },
          {
            "width": 450,
            "height": 450,
            "url": "https://media-cdn.tripadvisor.com/media/photo-s/1a/f6/ed/ce/default-avatar-2020-7.jpg"
          },
          {
            "width": 550,
            "height": 550,
            "url": "https://media-cdn.tripadvisor.com/media/photo-p/1a/f6/ed/ce/default-avatar-2020-7.jpg"
          },
          {
            "width": 1024,
            "height": 1024,
            "url": "https://media-cdn.tripadvisor.com/media/photo-w/1a/f6/ed/ce/default-avatar-2020-7.jpg"
          },
          {
            "width": 1200,
            "height": 1200,
            "url": "https://media-cdn.tripadvisor.com/media/photo-o/1a/f6/ed/ce/default-avatar-2020-7.jpg"
          }
        ]
      }
    },
    "username": "superReviewer"
  }
]

```