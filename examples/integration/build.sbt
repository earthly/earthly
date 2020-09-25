lazy val scalatest = "org.scalatest" %% "scalatest" % "3.0.5"

scalaVersion := "2.12.1"
name := "scala-example"
organization := "earthly.dev"
version := "1.0"

libraryDependencies ++= Seq(
  "org.tpolecat" %% "doobie-core"      % "0.9.0",
  "org.tpolecat" %% "doobie-postgres"  % "0.9.0",         
  "org.tpolecat" %% "doobie-scalatest" % "0.9.0" % "test" 
)

lazy val root = (project in file("."))
  .configs(IntegrationTest)
  .settings(
    Defaults.itSettings,
    libraryDependencies += scalatest % "it,test"
  )