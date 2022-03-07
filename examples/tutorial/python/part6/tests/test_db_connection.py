import unittest
import psycopg2

class MyIntegrationTests(unittest.TestCase):

    def test_db_connection_active(self):
        connection = psycopg2.connect(
            host="postgres",
            database="test_db",
            user="earthly",
            password="password")
        
        self.assertEqual(connection.closed, 0)

if __name__ == '__main__':
    unittest.main()
