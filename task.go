package main

type Task struct {
	interval int
	elapsed  int
	running  bool
	callback func()
}

func NewTask(interval int, callback func()) *Task {
	return &Task{
		interval: interval,
		elapsed:  0,
		callback: callback,
		running:  true,
	}
}

func (t *Task) Update() {
	t.elapsed++
	t.elapsed = (t.elapsed) % t.interval
	if t.elapsed == 0 {
		t.callback()
	}
}

func (t *Task) Pause() {
	t.running = false
}

func (t *Task) Resume() {
	t.running = true
}

func (t *Task) Reset() {
	t.elapsed = 0
}
