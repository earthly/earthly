import org.scalatest.FlatSpec
import doobie._
import doobie.implicits._
import cats.effect.IO
import scala.concurrent.ExecutionContext

class DatabaseIntegrationTest extends FlatSpec {
  implicit val cs = IO.contextShift(ExecutionContext.global)

  val xa = Transactor.fromDriverManager[IO](
    "org.postgresql.Driver", 
    "jdbc:postgresql://localhost:5432/iso3166", 
    "postgres",
    "postgres"
  )

  "A table" should "have country data" in {
    val dal = new DataAccessLayer()
    assert(dal.countries(5).transact(xa).unsafeRunSync.size == 5)
  }
}

