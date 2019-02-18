TalksTimer
=========

TalksTimer is a project that helps by timing conference/meetup talks by offering a 2 part timer. You can see this as a timer screen that you can remote control and send messages to.  
It starts by creating a session, this session gives you a key and the ability to change the time. When you enter this key on a 2nd screen (for the presenter) it will show you the timer.  
TalksTimer counts down on the server side, so it will continue should the controller/client disconnect.

## How to use
### Hosted
I host an instance at [https://talkstimer.com](https://talkstimer.com). It is free to use for everyone! If you happen to want to use this at a big (paid) conference I kindly ask you to email me, I pay all server costs out of my own pocket.

### Locally
1) Get the latest [Docker Compose](https://docs.docker.com/compose/install/)
2) Clone this repo
3) Run `docker-compose up` in the root of this repo
4) Open `http://localhost:8081`

## To Do
- [x] Add a frontend
- [x] Add messages
- [ ] Add client server connection state monitoring
- [ ] Add color on messages
- [x] Add color times (eg red for 5 minutes left)
- [x] Make hosted on a distributed system
- [x] Make all options configurable
