import doobie._
import doobie.implicits._
import cats.effect.IO
import scala.concurrent.ExecutionContext

object Main extends App {
  val dal = new DataAccessLayer()
  val dv = new DataVersion()

  if(dv.version() > 1)
  {
    implicit val cs = IO.contextShift(ExecutionContext.global)
    val xa = Transactor.fromDriverManager[IO](
      "org.postgresql.Driver", 
      "jdbc:postgresql://localhost:5432/iso3166", 
      "postgres",
      "postgres"
    )

    val countries = dal.countries(5)
                       .transact(xa).unsafeRunSync
                       .toList.map(_.name).mkString(", ")

    println(s"The first 5 countries alphabetically are: $countries")
  }
}

class DataAccessLayer()
{
  case class Country(name : String)

  def countries(limit : Int): ConnectionIO[List[Country]] 
      = sql"select name from country"
          .query[Country]
          .stream
          .take(limit)
          .compile.toList
}

class DataVersion(){
  def version() : Int  = 7
}