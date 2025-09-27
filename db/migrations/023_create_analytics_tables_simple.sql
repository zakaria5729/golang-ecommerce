-- Create simplified analytics tables for API-only tracking

-- Analytics visitors table (simplified)
CREATE TABLE IF NOT EXISTS analytics_visitors (
    id SERIAL PRIMARY KEY,
    visitor_id VARCHAR(255) UNIQUE NOT NULL,
    first_visit TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_visit TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    total_visits INTEGER DEFAULT 1,
    total_page_views INTEGER DEFAULT 0,
    unique_pages INTEGER DEFAULT 0,
    country VARCHAR(100),
    city VARCHAR(100),
    region VARCHAR(100),
    timezone VARCHAR(50),
    is_bot BOOLEAN DEFAULT FALSE,
    bot_name VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Analytics page views table (simplified)
CREATE TABLE IF NOT EXISTS analytics_page_views (
    id SERIAL PRIMARY KEY,
    visitor_id VARCHAR(255) NOT NULL,
    session_id VARCHAR(255) NOT NULL,
    user_id INT DEFAULT NULL,
    path VARCHAR(500) NOT NULL,
    page_title VARCHAR(500),
    http_method VARCHAR(10) DEFAULT 'GET',
    status_code INTEGER DEFAULT 200,
    response_time INTEGER DEFAULT 0, -- in milliseconds
    ip_address INET,
    user_agent TEXT,
    referrer TEXT,
    country VARCHAR(100),
    city VARCHAR(100),
    region VARCHAR(100),
    timezone VARCHAR(50),
    is_bot BOOLEAN DEFAULT FALSE,
    bot_name VARCHAR(100),
    visited_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (visitor_id) REFERENCES analytics_visitors(visitor_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Analytics sessions table (simplified)
CREATE TABLE IF NOT EXISTS analytics_sessions (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(255) UNIQUE NOT NULL,
    visitor_id VARCHAR(255) NOT NULL,
    user_id INT DEFAULT NULL,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP WITH TIME ZONE,
    duration INTEGER DEFAULT 0, -- in seconds
    page_views INTEGER DEFAULT 0,
    is_bounce BOOLEAN DEFAULT TRUE,
    referrer TEXT,
    entry_page VARCHAR(500),
    exit_page VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_by INT DEFAULT NULL,
    deleted_by INT DEFAULT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (visitor_id) REFERENCES analytics_visitors(visitor_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_analytics_visitors_visitor_id ON analytics_visitors(visitor_id);
CREATE INDEX IF NOT EXISTS idx_analytics_visitors_last_visit ON analytics_visitors(last_visit);
CREATE INDEX IF NOT EXISTS idx_analytics_visitors_is_bot ON analytics_visitors(is_bot);

CREATE INDEX IF NOT EXISTS idx_analytics_page_views_visitor_id ON analytics_page_views(visitor_id);
CREATE INDEX IF NOT EXISTS idx_analytics_page_views_session_id ON analytics_page_views(session_id);
CREATE INDEX IF NOT EXISTS idx_analytics_page_views_path ON analytics_page_views(path);
CREATE INDEX IF NOT EXISTS idx_analytics_page_views_visited_at ON analytics_page_views(visited_at);
CREATE INDEX IF NOT EXISTS idx_analytics_page_views_is_bot ON analytics_page_views(is_bot);

CREATE INDEX IF NOT EXISTS idx_analytics_sessions_session_id ON analytics_sessions(session_id);
CREATE INDEX IF NOT EXISTS idx_analytics_sessions_visitor_id ON analytics_sessions(visitor_id);
CREATE INDEX IF NOT EXISTS idx_analytics_sessions_started_at ON analytics_sessions(started_at);
