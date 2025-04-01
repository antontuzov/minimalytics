# Minimalytics


A minimalist privacy-focused analytics tool built with Go and SQLite, following Standard Go Project Layout.


## Features

- 🚀 Lightweight SQLite storage
- 📊 Dashboard with multiple visualizations:
  - Daily traffic trends
  - Unique vs total visitors
  - Top pages
  - Referrer sources
  - Device distribution
  - Browser usage
- 🔒 Basic authentication for dashboard access
- 🕵️ IP anonymization for privacy compliance
- 🗑️ Automatic data retention policy
- ⚡ Rate-limited tracking endpoint
- 📈 JSON API for all statistics
- 📱 Mobile-responsive interface

## Installation

### Prerequisites
- Go 1.23.+
- SQLite3

### Quick Start

1. Clone the repository:
```bash
git clone https://github.com/antontuzov/minimalytics.git
cd minimalytics