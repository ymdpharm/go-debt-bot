package main

import (
    "log"
    "net/http"
    "os"
    "strconv"
    "fmt"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/line/line-bot-sdk-go/linebot"
    "github.com/gomodule/redigo/redis"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	

	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

    channelSecret := os.Getenv("CHANNEL_SECRET")
    channelAccessToken := os.Getenv("CHANNEL_ACCESS_TOKEN")

	router.POST("/hook", func(c *gin.Context) {
        bot, err := linebot.New(channelSecret, channelAccessToken,)
        if err != nil {
            fmt.Println(err)
            return
        }
        events, err := bot.ParseRequest(c.Request)

        fmt.Println(os.Getenv("REDISTOGO_URL"))
        
        conn, err := redis.DialURL(os.Getenv("REDISTOGO_URL"))
        if err != nil {
            fmt.Println(err)
            return
        }
        defer conn.Close()
		
        for _, event := range events {
            if event.Type == linebot.EventTypeMessage {
                switch message := event.Message.(type) {
                case *linebot.TextMessage:
                    ans, err := getRes(event.Source ,message.Text, conn)
                    if ans != "" {
                        bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(ans)).Do()
                    }
                    if err != nil {
                        fmt.Println(err)
                    }
                }
            }
        }
	})
	router.Run(":" + port)
}

func getRes(source *linebot.EventSource, message string, conn redis.Conn) (string, error) {
    if source.Type != "group" { 
        return "グループに入れてくれないと困るんだけど", nil 
    }

    sliced := strings.Split(message, " ")
    price, err := strconv.Atoi(sliced[0])
    if err != nil {
        switch sliced[0] {
        case "iam": return storeNewUser(conn, source, sliced[1])
        case "check": return checkPrice(conn, source) 
        case "reset": return resetPrice(conn, source)
        case "help": return "see: https://github.com/ymdpharm/go-debt-bot/blob/master/README.md", nil
        }
    } else {
        return storePrice(conn, source, price)
    }
    return "", nil
}

func storeNewUser(conn redis.Conn, source *linebot.EventSource, name string) (string, error) {
    isExists, err := redis.Bool(conn.Do("EXISTS", source.GroupID + "_name_" + source.UserID))
    if err != nil {
        return "", err
    }
    if isExists {
        _, err = conn.Do("SET", source.GroupID + "_name_" + source.UserID, name)
        if err != nil {
            return "", err
        }
        return "呼び方変えるわ", nil
    }
    _, err = conn.Do("SET", source.GroupID + "_name_" + source.UserID, name)
    ;conn.Do("SADD", source.GroupID + "_all", source.UserID)
    ;conn.Do("SET", source.GroupID + "_price_" + source.UserID, 0)
    if err != nil { 
        return "", err
    } else {
        return name + "ね, 覚えた", nil
    }
}

func checkPrice(conn redis.Conn, source *linebot.EventSource) (string, error) {
    isExists, err := redis.Bool(conn.Do("EXISTS", source.GroupID + "_all"))
    if err != nil {
        return "", err
    } else if !isExists {
        return "最初から check はきつい", nil
    }

    allUsers, err := redis.Strings(conn.Do("SMEMBERS", source.GroupID + "_all"))
    if err != nil {
        return "", err
    }

    userNum := len(allUsers)
    names := make([]string, userNum)
    prices := make([]int, userNum)
    var priceMin int

    for ind, elem := range allUsers {
        name, err := redis.String(conn.Do("GET", source.GroupID + "_name_" + elem))
        if err != nil {
            return "", err
        }

        price, err := redis.Int(conn.Do("GET", source.GroupID + "_price_" + elem))
        if err != nil {
            return "", err
        }

        names[ind] = name
        prices[ind] = price

        if priceMin == 0 || priceMin > price {
            priceMin = price
        }
    }

    ans := "今の貸し借りはね,"
    for i:= 0; i < userNum; i++ {
        ans += "\n"
        ans += names[i] + "に" + strconv.Itoa(prices[i] - priceMin) + "円" 
    }

    return ans, nil
}

func resetPrice(conn redis.Conn, source *linebot.EventSource) (string, error) {
    isExists, err := redis.Bool(conn.Do("EXISTS", source.GroupID + "_all"))
    if err != nil {
        return "", nil
    } else if !isExists {
        return "最初から reset はきつい", nil
    }

    allUsers, err := redis.Strings(conn.Do("SMEMBERS", source.GroupID + "_all"))
    if err != nil {
        return "", err
    }
    for _, elem := range allUsers {
        _, err := conn.Do("SET", source.GroupID + "_price_" + elem, 0)
        if err != nil {
            return "", err
        }
    }
    return "全部ゼロにした", nil
}

func storePrice(conn redis.Conn, source *linebot.EventSource, price int) (string, error) {
    isExists, err := redis.Bool(conn.Do("EXISTS", source.GroupID + "_name_" + source.UserID))
    if err != nil {
        return "", nil
    } else if !isExists {
        return "誰? iam ** って言ってからにして", nil
    }
    name, err := redis.String(conn.Do("GET", source.GroupID + "_name_" + source.UserID))
    if err != nil {
        return "", err
    } 

    prev, err := redis.Int(conn.Do("GET", source.GroupID + "_price_" + source.UserID))
    if err != nil {
        return "", err
    }

    if _, err := conn.Do("SET", source.GroupID + "_price_" + source.UserID, prev + price); err != nil {
        return "", err
    } else { return name + "に" + strconv.Itoa(price) + "円ね", nil }
}
