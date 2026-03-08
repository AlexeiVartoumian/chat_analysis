-- ============================================================
-- 1. TOP KEYWORDS PER CHAT
-- ============================================================
SELECT 
    c.title,
    cc.word,
    cc.frequency,
    cc.tf_idf_score
FROM chat_concepts cc
JOIN chats c ON cc.chat_id = c.chat_id
WHERE cc.role = 'conversation'
AND cc.chat_id = 'insert-uuid-of-entire-chat-here'
ORDER BY cc.tf_idf_score DESC
LIMIT 10;


-- ============================================================
-- 2. CHATS THAT SHARE KEYWORDS
-- ============================================================
SELECT 
    a.chat_id        AS chat_a,
    b.chat_id        AS chat_b,
    a.word,
    a.tf_idf_score + b.tf_idf_score AS edge_weight
FROM chat_concepts a
JOIN chat_concepts b 
    ON  a.word    = b.word
    AND a.chat_id < b.chat_id        
    AND a.role    = b.role
WHERE a.role = 'full'
ORDER BY edge_weight DESC
LIMIT 20;


-- ============================================================
-- 3. HUMAN VS ASSISTANT KEYWORD COMPARISON
-- ============================================================
SELECT
    h.word,
    h.frequency      AS human_frequency,
    h.tf_idf_score   AS human_tfidf,
    a.frequency      AS assistant_frequency,
    a.tf_idf_score   AS assistant_tfidf,
    h.tf_idf_score - a.tf_idf_score AS human_dominance
FROM chat_concepts h
JOIN chat_concepts a 
    ON  h.word    = a.word
    AND h.chat_id = a.chat_id
WHERE h.role = 'human'
AND   a.role = 'assistant'
AND   h.chat_id = 'insert-uuid-of-entire-chat-here'
ORDER BY human_dominance DESC;


-- ============================================================
-- 4. KEYWORD FREQUENCY ACROSS ALL CHATS
-- ============================================================
SELECT
    word,
    COUNT(DISTINCT chat_id)  AS chat_count,
    SUM(frequency)           AS total_frequency,
    AVG(tf_idf_score)        AS avg_tfidf,
    MAX(tf_idf_score)        AS peak_tfidf
FROM chat_concepts
WHERE role = 'full'
GROUP BY word
ORDER BY avg_tfidf DESC
LIMIT 20;
