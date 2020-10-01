import org.scalatest.FlatSpec

class DataVersionSpec extends FlatSpec {

  val dv = new DataVersion()
  "Data Version " should " be positive" in {
    assert(dv.version > 0)
  }
}