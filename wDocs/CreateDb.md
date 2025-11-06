# 🗄️ Auth Service Database Setup Guide

This guide explains how to create and initialize the **Auth Service** PostgreSQL database for the Social Network backend.

---

## ⚙️ 1. Prerequisites

Make sure you have PostgreSQL installed and running.

### 🧩 Local installation

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib -y
```

```bash
sudo -u postgres psql


-- Create a new database for the Auth service
CREATE DATABASE auth;
```