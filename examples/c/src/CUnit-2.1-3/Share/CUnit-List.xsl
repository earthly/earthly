<?xml version='1.0'?>
<xsl:stylesheet	version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">

	<xsl:template match="CUNIT_TEST_LIST_REPORT">
		<html>
			<head>
				<title> CUnit - Suite and Test Case Organization in Test Registry </title>
			</head>

			<body bgcolor="#e0e0f0">
				<xsl:apply-templates/>
			</body>
		</html>
	</xsl:template>

	<xsl:template match="CUNIT_HEADER">
		<div align="center">
			<h3>
				<b> CUnit - A Unit testing framework for C </b> <br/>
				<a href="http://cunit.sourceforge.net/"> http://cunit.sourceforge.net/ </a>
			</h3>
		</div>
	</xsl:template>

	<xsl:template match="CUNIT_LIST_TOTAL_SUMMARY">
		<p/>
		<table align="center" width="50%">
			<xsl:apply-templates/>
		</table>
	</xsl:template>

	<xsl:template match="CUNIT_LIST_TOTAL_SUMMARY_RECORD">
		<tr>
			<td bgcolor="#f0f0e0" width="70%">
				<xsl:value-of select="CUNIT_LIST_TOTAL_SUMMARY_RECORD_TEXT" />
			</td>
			<td bgcolor="#f0e0e0" align="center">
				<xsl:value-of select="CUNIT_LIST_TOTAL_SUMMARY_RECORD_VALUE" />
			</td>
		</tr>
	</xsl:template>

	<xsl:template match="CUNIT_ALL_TEST_LISTING">
		<p/>
		<div align="center">
			<h2> Listing of Registered Suites &amp; Tests </h2>
		</div>
		<table align="center" width="90%">
			<tr bgcolor="#00ccff">
				<td colspan="8"> </td>
			</tr>
			<tr>
				<td width="44%" colspan="2" bgcolor="#f0f0e0"> </td>
				<td width="14%" bgcolor="#f0f0e0" align="center"> <b> Initialize Function? </b> </td>
				<td width="14%" bgcolor="#f0f0e0" align="center"> <b> Cleanup Function? </b> </td>
				<td width="14%" bgcolor="#f0f0e0" align="center"> <b> Test Count </b> </td>
				<td width="14%" bgcolor="#f0f0e0" align="center"> <b> Active? </b> </td>
			</tr>
			<xsl:apply-templates/>
		</table>
	</xsl:template>

	<xsl:template match="CUNIT_ALL_TEST_LISTING_SUITE">
		<xsl:apply-templates/>
	</xsl:template>

	<xsl:template match="CUNIT_ALL_TEST_LISTING_SUITE_DEFINITION">
		<tr bgcolor="#00ccff">
			<td colspan="8"> </td>
		</tr>
		<tr>
			<td bgcolor="#f0e0f0" align="left"> Suite </td>
			<td bgcolor="#f0e0f0" align="left"> <b> <xsl:value-of select="SUITE_NAME" /> </b> </td>

			<td bgcolor="#f0e0f0" align="center"> <xsl:value-of select="INITIALIZE_VALUE" /> </td>
			<td bgcolor="#f0e0f0" align="center"> <xsl:value-of select="CLEANUP_VALUE" /> </td>
			<td bgcolor="#f0e0f0" align="center"> <xsl:value-of select="TEST_COUNT_VALUE" /> </td>
			<td bgcolor="#f0e0f0" align="center"> <xsl:value-of select="ACTIVE_VALUE" /> </td>
		</tr>
		<tr bgcolor="#00ccff">
			<td colspan="8"> </td>
		</tr>
	</xsl:template>

	<xsl:template match="CUNIT_ALL_TEST_LISTING_SUITE_TESTS">
		<xsl:apply-templates/>
	</xsl:template>

	<xsl:template match="TEST_CASE_DEFINITION">
		<tr>
			<td align="center" bgcolor="#e0f0d0"> Test </td>
			<td align="left" colspan="4" bgcolor="#e0e0d0">
				<xsl:value-of select="TEST_CASE_NAME" />
			</td>
			<td align="center" bgcolor="#e0e0d0">
				<xsl:value-of select="TEST_ACTIVE_VALUE" />
			</td>
		</tr>
	</xsl:template>

	<xsl:template match="CUNIT_ALL_TEST_LISTING_GROUP">
		<xsl:apply-templates/>
	</xsl:template>

	<xsl:template match="CUNIT_ALL_TEST_LISTING_GROUP_DEFINITION">
		<tr>
			<td width="10%" bgcolor="#f0e0f0"> Suite </td>
			<td width="20%" bgcolor="#e0f0f0" >
				<b> <xsl:value-of select="GROUP_NAME" /> </b>
			</td>

			<td width="15%" bgcolor="#f0e0f0"> Initialize Function? </td>
			<td width="5%" bgcolor="#e0f0f0">
				<xsl:value-of select="INITIALIZE_VALUE" />
			</td>

			<td width="15%" bgcolor="#f0e0f0"> Cleanup Function? </td>
			<td width="5%" bgcolor="#e0f0f0">
				<xsl:value-of select="CLEANUP_VALUE" />
			</td>

			<td width="10%" bgcolor="#f0e0f0"> Test Count </td>
			<td width="5%" bgcolor="#e0f0f0">
				<xsl:value-of select="TEST_COUNT_VALUE" />
			</td>
		</tr>
	</xsl:template>

	<xsl:template match="CUNIT_ALL_TEST_LISTING_GROUP_TESTS">
		<tr>
			<td align="center" bgcolor="#e0f0d0"> Test Cases </td>
			<td align="left" colspan="7" bgcolor="#e0e0d0">
				<xsl:for-each select="TEST_CASE_NAME">
					<xsl:apply-templates/> <br />
				</xsl:for-each>
			</td>
		</tr>
		<tr bgcolor="#00ccff">
			<td colspan="8"> </td>
		</tr>
	</xsl:template>

	<xsl:template match="CUNIT_FOOTER">
		<p/>
		<hr align="center" width="90%" color="maroon" />
		<h5 align="center">
			<xsl:apply-templates/>
		</h5>
	</xsl:template>

</xsl:stylesheet>
