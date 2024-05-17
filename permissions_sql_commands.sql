-- Set the activated field for alice@example.com to true.
UPDATE members SET activated = true WHERE email = 'alice@example.com';


-- Give all members the 'movies:read' permission
INSERT INTO members_permissions
SELECT id, (SELECT id FROM permissions WHERE code = 'movies:read') FROM members;


-- Give faith@example.com the 'movies:write' permission
INSERT INTO members_permissions
    VALUES (
    (SELECT id FROM members WHERE email = 'faith@example.com'),
    (SELECT id FROM permissions WHERE code = 'movies:write')
);


-- List all activated members and their permissions.
SELECT email, array_agg(permissions.code) as permissions
FROM permissions
INNER JOIN members_permissions ON members_permissions.permission_id = permissions.id
INNER JOIN members ON members_permissions.member_id = members.id
WHERE members.activated = true
GROUP BY email;