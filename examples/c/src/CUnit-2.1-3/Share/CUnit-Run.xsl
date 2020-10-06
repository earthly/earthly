<?xml version='1.0'?>
<xsl:stylesheet
	version="1.0"
	xmlns:xsl="http://www.w3.org/1999/XSL/Transform">

	<xsl:template match="CUNIT_TEST_RUN_REPORT">
		<html>
			<head>
				<title> CUnit - Automated Test Run Summary Report </title>
			</head>

			<body bgcolor="#e0e0f0">
				<xsl:apply-templates/>
			</body>
		</html>
	</xsl:template>

	<xsl:template match="CUNIT_HEADER">
		<div align="center">
			<h3>
				<b> CUnit - A Unit testing framework for C. </b> <br/>
				<a href="http://cunit.sourceforge.net/"> http://cunit.sourceforge.net/ </a>
			</h3>
		</div>
	</xsl:template>

	<xsl:template match="CUNIT_RESULT_LISTING">
		<p/>
		<div align="center">
			<h2> Automated Test Run Results </h2>
		</div>
		<table cols="4" width="90%" align="center">
			<tr>
				<td width="25%"> </td>
				<td width="25%"> </td>
				<td width="25%"> </td>
				<td width="25%"> </td>
			</tr>
			<xsl:apply-templates/>
		</table>
	</xsl:template>

	<xsl:template match="CUNIT_RUN_SUITE">
		<xsl:apply-templates/>
	</xsl:template>

	<xsl:template match="SUITE_NAME">
	</xsl:template>

	<xsl:template match="CUNIT_RUN_SUITE_SUCCESS">
		<tr bgcolor="#f0e0f0">
			<td colspan="4">
				Running Suite <xsl:value-of select="SUITE_NAME"/>
			</td>
		</tr>
		<xsl:apply-templates/>
	</xsl:template>

	<xsl:template match="CUNIT_RUN_GROUP">
		<xsl:apply-templates/>
	</xsl:template>

	<xsl:template match="CUNIT_RUN_GROUP_SUCCESS">
		<tr bgcolor="#f0e0f0">
			<td colspan="4">
				Running Group <xsl:apply-templates/>
			</td>
		</tr>
	</xsl:template>

	<xsl:template match="CUNIT_RUN_TEST_RECORD">
		<xsl:apply-templates/>
	</xsl:template>

	<xsl:template match="CUNIT_RUN_TEST_SUCCESS">
		<tr bgcolor="#e0f0d0">
			<td> </td>
			<td colspan="2">
				Running test <xsl:apply-templates/>...
			</td>
			<td bgcolor="#50ff50"> Passed </td>
		</tr>
	</xsl:template>

	<xsl:template match="CUNIT_RUN_TEST_FAILURE">
		<tr bgcolor="#e0f0d0">
			<td> </td>
			<td colspan="2">
				Running test <xsl:value-of select="TEST_NAME"/>...
			</td>
			<td bgcolor="#ff5050"> Failed </td>
		</tr>

		<tr>
			<td colspan="4" bgcolor="#ff9090">
				<table width="100%">
					<tr>
						<th width="15%"> File Name </th>
						<td width="50%" bgcolor="#e0eee0">
							<xsl:value-of select="FILE_NAME"/>
						</td>
						<th width="20%"> Line Number </th>
						<td width="10%" bgcolor="#e0eee0">
							<xsl:value-of select="LINE_NUMBER"/>
						</td>
					</tr>
					<tr>
						<th width="15%"> Condition </th>
						<td colspan="3" width="85%" bgcolor="#e0eee0">
							<xsl:value-of select="CONDITION"/>
						</td>
					</tr>
				</table>
			</td>
		</tr>
	</xsl:template>

	<xsl:template match="CUNIT_RUN_SUITE_FAILURE">
		<tr>
			<td colspan="3" bgcolor="#f0b0f0">
				Running Suite <xsl:value-of select="SUITE_NAME"/>...
			</td>
			<td bgcolor="#ff7070">
				<xsl:value-of select="FAILURE_REASON"/>
			</td>
		</tr>
	</xsl:template>

	<xsl:template match="CUNIT_RUN_GROUP_FAILURE">
		<tr>
			<td colspan="3" bgcolor="#f0b0f0">
				Running Group <xsl:value-of select="GROUP_NAME"/>...
			</td>
			<td bgcolor="#ff7070">
				<xsl:value-of select="FAILURE_REASON"/>
			</td>
		</tr>
	</xsl:template>

	<xsl:template match="CUNIT_RUN_SUMMARY">
		<p/>
		<table width="90%" rows="5" align="center">
			<tr align="center" bgcolor="skyblue">
				<th colspan="6"> Cumulative Summary for Run </th>
			</tr>
			<tr>
				<th width="15%" bgcolor="#ffffc0" align="center"> Type </th>
				<th width="17%" bgcolor="#ffffc0" align="center"> Total </th>
				<th width="17%" bgcolor="#ffffc0" align="center"> Run </th>
				<th width="17%" bgcolor="#ffffc0" align="center"> Succeeded </th>
				<th width="17%" bgcolor="#ffffc0" align="center"> Failed </th>
				<th width="17%" bgcolor="#ffffc0" align="center"> Inactive </th>
			</tr>
			<xsl:for-each select="CUNIT_RUN_SUMMARY_RECORD">
				<tr align="center" bgcolor="lightgreen">
					<td> <xsl:value-of select="TYPE" /> </td>
					<td> <xsl:value-of select="TOTAL" /> </td>
					<td> <xsl:value-of select="RUN" /> </td>
					<td> <xsl:value-of select="SUCCEEDED" /> </td>
					<td> <xsl:value-of select="FAILED" /> </td>
					<td> <xsl:value-of select="INACTIVE" /> </td>
				</tr>
			</xsl:for-each>
		</table>
	</xsl:template>

	<xsl:template match="CUNIT_FOOTER">
		<p/>
		<hr align="center" width="90%" color="maroon" />
		<h5 align="center"> <xsl:apply-templates/> </h5>
	</xsl:template>

</xsl:stylesheet>
