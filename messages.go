package main

type Message interface{}

type Messages []Message

type MSG_NIL struct {}

type MSG_REMOVE struct {}
type MSG_CLEAR struct { x, y int }