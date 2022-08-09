DROP TABLE IF EXISTS task
CREATE TABLE task (
    id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    detail TEXT NOT NULL,
    assignee VARCHAR(100) NOT NULL,
    deadline TIMESTAMP NULL,
    status ENUM('undone', 'done') DEFAULT 'undone',
    PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;