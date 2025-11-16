# 🧠 My System Design Bootcamp (Self-Inflicted Edition)

*aka: “I Don’t Just Read About Systems, I Accidentally Break Them.”*

Welcome to my daily(ish) 90-minute chaos sessions, where I attempt to learn System Design properly — not by reading blogs (sorry to blog lovers :(), but by building stuff, breaking stuff, fixing stuff, and pretending this was the plan all along.

![Coding GIF](https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExNTFkN2U2ZDI2NjQ5YjJkZjE2YjBlY2Q5NTcwM2M4Y2M2OTQ2ZDAwZCZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/13HgwGsXF0aiGY/giphy.gif)

---

### 🚀 What’s Happening Here?

This repo contains my day-by-day system design progress, where I:

1.  Write hundreds of lines of Go code 👨‍💻
2.  Break things on purpose 💥
3.  Break things *not* on purpose 😭
4.  Then measure *why* it broke 📊
5.  Then pretend I’m learning from it 🤔
6.  Then actually learn something 💡
7.  Then break it again, because repetition builds muscle memory 💪

### 🤔 The Learning Philosophy

You know those tutorials that say: “*First, let’s read about Redis.*”

Lol no.

I build the cache first, break it, benchmark it, stress it, then cry, then fix it, then I look at Redis and say:
> “Ah. *That’s* why you do it like that.”

Every concept → implemented before understood.
Understanding comes from the suffering.
This is the way.

![alt text](https://media.giphy.com/media/qs6ev2pm8g9dS/giphy.gif)

---

### 📚 Project Index & Key Learnings

*A guide to the chaos. Each entry links to the code and my distilled findings.*

#### **[Week 1: In-Memory LRU Cache](./week-01-caching)** 🧠
*   **Mission:** Build a thread-safe, performant LRU cache from scratch to understand the fundamentals of caching systems.
*   **Key Finding:** High contention can be **faster** than low contention if it avoids memory allocation. A cache *hit* is a cheap CPU operation; a cache *miss* is an expensive memory operation. The cost of memory allocation can be far greater than the cost of lock contention.

#### **[Week 2: TCP Load Balancer](./week-02-load-balancer)** ⚖️
*   **Mission:** *(Coming Soon)*
*   *Key Finding:* *(Coming Soon)*

#### **[Week 3: Message Queue](./week-03-message-queue)** 📬
*   **Mission:** *(Coming Soon)*
*   *Key Finding:* *(Coming Soon)*

---

### 🛠️ Tech Stack

Because apparently I enjoy pain, but not *that* much:

*   **Go** (the official language of “I swear this goroutine leak wasn’t my fault”)
*   **Docker** (so I can break things consistently across machines)
*   **AWS Free Tier** (aka: “I hope I don't accidentally get billed for a small moon”)
*   **SQL** (starting fresh anyway because why not suffer twice)

### 📅 Time Commitment

**90 mins/day × 6 days/week** = Just Enough Time To Regret My Life Choices

---

### 📦 What You’ll Find in This Repo:

*   Daily folders/logs 📂
*   Code that works ✅
*   Code that *should* work... 🤨
*   Code that definitely should NOT work ❌
*   Notes, benchmarks, diagrams, regrets 📝
*   Occasional forehead imprints on the keyboard that I forgot to delete

### 🧪 Expected Side Effects:

*   Sudden understanding of bottlenecks
*   A compulsive need to over-optimize
*   Random goroutines multiplying in the background
*   Calling everything a “distributed system” even if it's just two Go files
*   A dangerous amount of confidence

---

### 🫠 How to Use This Repo

1.  Clone it.
2.  Open anything.
3.  Realize I have no idea what I’m doing.
4.  Watch as I somehow figure it out anyway.
5.  Become inspired. ✨
6.  Start your own chaos.
7.  We rise together. 🫡

### 🙌 Let’s Go Build Something.

Or break something.
<br>Either way — progress.

🔨