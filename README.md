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

[]

```