// Copyright (C) 2014 Zenoss, Inc
//
// redis-mon is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 2 of the License, or
// (at your option) any later version.
//
// redis-mon is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Foobar. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisStats struct {
	connString string
}

type RedisListLength struct {
	RedisStats
	listName string
}

var tags map[string]string

func init() {
	tags = make(map[string]string)
	if h, err := os.Hostname(); err != nil {
		panic(err)
	} else {
		tags["hostname"] = h
	}
}

func NewRedisStats(connString string) *RedisStats {
	s := &RedisStats{
		connString: connString,
	}
	return s
}

func NewRedisListLength(connString, listName string) *RedisListLength {
	return &RedisListLength{
		RedisStats{
			connString: connString,
		},
		listName,
	}
}

func (s *RedisListLength) Read() (stats []Sample, err error) {
	c, err := redis.Dial("tcp", s.connString)
	if err != nil {
		return nil, err
	}
	val, err := redis.Int(c.Do("LLEN", s.listName))
	if err != nil {
		return nil, err
	}
	now := time.Now().Unix()
	stats = []Sample{
		Sample{
			Metric:    fmt.Sprintf("redis.%sLength", s.listName),
			Value:     fmt.Sprintf("%d", val),
			Timestamp: now,
			Tags:      tags,
		},
	}
	return stats, nil
}

func (s *RedisStats) Read() (stats []Sample, err error) {
	c, err := redis.Dial("tcp", s.connString)
	if err != nil {
		return nil, err
	}
	val, err := redis.String(c.Do("INFO"))
	if err != nil {
		return nil, err
	}
	now := time.Now().Unix()

	stats = make([]Sample, 0)
	str := string(val)
	for _, line := range strings.Split(str, "\n") {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}
		_, err := strconv.ParseFloat(parts[1], 32)
		if err != nil {
			continue
		}
		stats = append(stats, Sample{
			Metric:    "redis." + parts[0],
			Value:     parts[1],
			Timestamp: now,
			Tags:      tags,
		})
	}
	return stats, nil
}
