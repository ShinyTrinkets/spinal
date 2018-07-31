---
id: pomodoro-timer
trinkets: false
log: false
db: false
---

# Pomodoro 🍅 timer ⏰

The [Pomodoro Technique](https://en.wikipedia.org/wiki/Pomodoro_Technique) is a time management method that uses a timer to break down work into intervals, traditionally 25 minutes, separated by short breaks.

```js
// Take a break
trigger('timer', '25,55 9-18 * * 1-5', action_break)
// Back to work
trigger('timer', '00,30 9-18 * * 1-5', action_work)
```

And here are the actions:

```js
function action_break () {
  console.warn('Take a break 💤')
}
function action_work () {
  console.warn('Back to work 💪')
}
```

## Good bye 🕰
