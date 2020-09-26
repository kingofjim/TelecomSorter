package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yl2chen/cidranger"
	"log"
	"net"
	"runtime"
	"strings"
	"time"
)

type ranger struct {
	id      int16
	country string
	state   string
	telecom string
	ranger  cidranger.Ranger
}

func Ranger(c *gin.Context) {
	start := time.Now()
	ip := c.Query("ip")
	if ip == "" {
		c.JSON(200, gin.H{
			"query": false,
		})
		return
	} else {
		//PrintMemUsage()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		outChan := make(chan ranger)
		for _, ranger := range ranger_army {
			go rangerSearch(ctx, outChan, ranger, net.ParseIP(ip))
		}
		//cancel()
		found := false
		var result ranger
		LOOP:
			for i := 0; i < len(ranger_army); i++ {
				fmt.Printf("Index: %v\n", i)
				select {
				case ranger := <-outChan:
					//fmt.Println("finished:", ranger)
					result = ranger
					cancel()
					<-ctx.Done()
					found = true
					break LOOP
				default:
				}
			}
		if !found {
			fmt.Printf("%v - Not Found\n", ip)
			c.JSON(200, gin.H{
				"query": false,
			})
		} else {
			fmt.Printf("%v - %v\n", ip, result)
			c.JSON(200, gin.H{
				"query":   true,
				"ip":      ip,
				"country": result.country,
				"state":   result.state,
				"telecom": result.telecom,
			})
		}
	}
	elapsed := time.Since(start)
	log.Printf("Mapping took %s", elapsed)

	//PrintMemUsage()

}

func buildArmy() []ranger {
	var tempRanger []ranger
	db, err := sql.Open("sqlite3", "ip_library_converter/sqlite")
	checkErr(err)
	rows, err := db.Query("SELECT id, country, state, telecom, ip FROM ip_table")
	checkErr(err)
	//fmt.Print(rows)

	for rows.Next() {
		var ip, country, state, telecom string
		var id int16
		err = rows.Scan(&id, &country, &state, &telecom, &ip)
		checkErr(err)
		//fmt.Println(ip)
		ip_slice := strings.Split(ip, ",")
		ranger := ranger{id: id, country: "China", state: state, telecom: telecom, ranger: cidranger.NewPCTrieRanger()}
		for _, ip := range ip_slice {
			_, network, _ := net.ParseCIDR(ip)
			ranger.ranger.Insert(cidranger.NewBasicRangerEntry(*network))
		}
		tempRanger = append(tempRanger, ranger)
	}
	return tempRanger
}

func rangerSearch(ctx context.Context, outChan chan<- ranger, ranger ranger, ip net.IP) {
	if ip != nil {
		check, err := ranger.ranger.Contains(ip)
		//fmt.Printf("%v %v\n", ranger, check)
		if checkErrPanic(err) {
			if check {
				outChan <- ranger
			}
		}
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkErrPanic(err error) bool {
	if err != nil {
		log.Panic(err)
		return false
	}
	return true
}
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
	fmt.Print("\n")
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
