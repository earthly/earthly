package postgresclient;

import org.joda.time.LocalTime;
import java.sql.Connection;
import java.sql.DriverManager;


public class PostgreSQLJDBC {
   public static void main(String args[]) {
      Connection c = null;
      try {
         Class.forName("org.postgresql.Driver");
         c = DriverManager
            .getConnection("jdbc:postgresql://postgres:5432/test_db",
            "earthly", "password");
      } catch (Exception e) {
         e.printStackTrace();
         System.err.println(e.getClass().getName()+": "+e.getMessage());
         System.exit(0);
      }
      System.out.println("Opened database successfully!");
   }
}
