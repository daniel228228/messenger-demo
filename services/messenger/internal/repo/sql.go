package repo

var (
	dialogQueryOffset = `WITH B AS
  (SELECT p.from_user_id,
          p.to_user_id,
          m.*
   FROM message_peer_user_to_user p
   JOIN message m ON m.id = p.message_id
   WHERE p.from_user_id = $1
     OR p.to_user_id = $1 ),
     R AS
  (SELECT b1.*
   FROM B b1
   JOIN
     (SELECT from_user_id,
             to_user_id,
             max(timestamp) AS max_timestamp
      FROM B
      GROUP BY from_user_id,
               to_user_id) b2 ON b2.max_timestamp = b1.timestamp
   AND b2.from_user_id = b1.from_user_id
   AND b2.to_user_id = b1.to_user_id)
SELECT *
FROM R R1
WHERE NOT EXISTS
    (SELECT *
     FROM R R2
     WHERE R1.from_user_id = R2.to_user_id
       AND R1.to_user_id = R2.from_user_id
       AND R1.timestamp < R2.timestamp )
  AND timestamp < $2
ORDER BY timestamp DESC
LIMIT $3`

	dialogQueryTotal = `WITH R AS
  (SELECT max(p.message_id) AS max_message_id,
          p.from_user_id,
          p.to_user_id
   FROM message_peer_user_to_user p
   WHERE p.from_user_id = $1
     OR p.to_user_id = $1
   GROUP BY p.from_user_id,
            p.to_user_id)
SELECT COUNT(*)
FROM R r1
WHERE NOT EXISTS
    (SELECT *
     FROM R R2
     WHERE R1.from_user_id = R2.to_user_id
       AND R1.to_user_id = R2.from_user_id
       AND R1.max_message_id < R2.max_message_id )`

	msgQuery = `SELECT p.from_user_id,
       m.*
FROM message_peer_user_to_user p
JOIN message m ON m.id = p.message_id
WHERE ((p.from_user_id = $1
         AND p.to_user_id = $2)
        OR (p.from_user_id = $2
            AND p.to_user_id = $1))
ORDER BY timestamp DESC
LIMIT $3`

	msgQueryTimestamp = `SELECT m.timestamp
FROM message m
JOIN message_peer_user_to_user p ON p.message_id = m.id
AND ((p.from_user_id = $1
      AND p.to_user_id = $2)
     OR (p.from_user_id = $2
         AND p.to_user_id = $1))
WHERE m.id = $3`

	msgQueryOffset = `SELECT p.from_user_id,
       m.*
FROM message_peer_user_to_user p
JOIN message m ON m.id = p.message_id
AND (timestamp, id) < ($1,
                       $2)
WHERE ((p.from_user_id = $3
         AND p.to_user_id = $4)
        OR (p.from_user_id = $4
            AND p.to_user_id = $3))
ORDER BY timestamp DESC
LIMIT $5`

	msgQueryTotal = `SELECT COUNT(*)
FROM message_peer_user_to_user p
JOIN message m ON m.id = p.message_id
WHERE ((p.from_user_id = $1
         AND p.to_user_id = $2)
        OR (p.from_user_id = $2
            AND p.to_user_id = $1))`

	msgReadQuery = `
  INSERT INTO message_read_user_to_user (
      from_user_id,
      to_user_id,
      last_read_message_id
    )
  SELECT $1,
    $2,
    $3
  WHERE EXISTS (
      SELECT *
      FROM message_peer_user_to_user p
      WHERE (
          p.message_id = $3
          AND p.from_user_id = $2
          AND p.to_user_id = $1
        )
    ) ON CONFLICT (from_user_id, to_user_id) DO
  UPDATE
  SET last_read_message_id = $3
  WHERE EXISTS (
      SELECT *
      FROM messages m1
        JOIN messages m2 ON m1.id = $3
        AND m2.id = message_read_user_to_user.last_read_message_id
        AND m1.timestamp > m2.timestamp
    )`

	msgUnreadQueryCount = `SELECT COUNT(*)
FROM message_peer_user_to_user p
JOIN message m ON m.id = p.message_id
WHERE p.from_user_id = $1
  AND p.to_user_id = $2`

	msgUnreadQueryCountOffset = `SELECT COUNT(*)
FROM message_peer_user_to_user p
JOIN message m ON m.id = p.message_id
AND ((timestamp, id) > ($1,
                       $2) )
WHERE p.from_user_id = $3
  AND p.to_user_id = $4`

	isUnreadDialogQuery = `WITH B AS (
SELECT
	m.timestamp,
	m.id
FROM
	message_read_user_to_user r
JOIN message m ON
	m.id = r.last_read_message_id
WHERE
	r.from_user_id = $1
	AND r.to_user_id = $2 )
SELECT
	1 AS is_unread
WHERE
	EXISTS (
	SELECT
		*
	FROM
		message_peer_user_to_user p
	JOIN message m ON
		m.id = p.message_id
		AND ( (timestamp, id) > ( SELECT * FROM B )
		OR NOT EXISTS ( SELECT * FROM B) )
	WHERE
		p.from_user_id = $2
		AND p.to_user_id = $1
)`
)
