package controller

import (
    "net/http"
    "strconv"
)

func ParseUintQueryParam(r *http.Request, key string) (uint, error) {
    val := r.URL.Query().Get(key)
    parsed, err := strconv.ParseUint(val, 10, 64)
    return uint(parsed), err
}
