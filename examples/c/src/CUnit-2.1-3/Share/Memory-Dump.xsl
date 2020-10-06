<?xml version='1.0'?>
<xsl:stylesheet
	version="1.0"
	xmlns:xsl="http://www.w3.org/1999/XSL/Transform">

	<xsl:template match="MEMORY_DUMP_REPORT">
		<html>
			<head>
				<title> CUnit Memory Debugger Dumper - All Allocation/Deallocation Report. </title>
			</head>

			<body bgcolor="e0e0f0">
				<xsl:apply-templates/>
			</body>
		</html>
	</xsl:template>

	<xsl:template match="MD_HEADER">
		<div align="center">
			<h3>
				<b> CUnit - A Unit testing framework for C. </b> <br/>
				<a href="http://cunit.sourceforge.net/"> http://cunit.sourceforge.net/ </a>
			</h3>
		</div>
	</xsl:template>

	<xsl:template match="MD_RUN_LISTING">
		<div align="center">
			<h2>
				CUnit Memory Debugger Report <br/>
				Memory Allocation/Deallocation Records
			</h2>
			<hr align="center" width="100%" color="maroon" />
		</div>
		<table cols="6" width="95%">
			<th width="10%" align="left">Pointer</th>
			<th width="35%" align="left">Allocation File</th>
			<th width="10%" align="left">Line #</th>
			<th width="35%" align="left">Deallocation File</th>
			<th width="10%" align="left">Line #</th>
			<th width="10%" align="left">Data Size</th>
			<xsl:apply-templates/>
		</table>
	</xsl:template>

	<xsl:template match="MD_RUN_RECORD">
		<xsl:param name="ptr" select="MD_POINTER"/>
		<xsl:for-each select="MD_EVENT_RECORD">
			<tr>
				<td> <xsl:value-of select="$ptr"/> </td>
				<td> <xsl:value-of select="MD_ALLOC_FILE"/> </td>
				<td> <xsl:value-of select="MD_ALLOC_LINE"/> </td>
				<td> <xsl:value-of select="MD_DEALLOC_FILE"/> </td>
				<td> <xsl:value-of select="MD_DEALLOC_LINE"/> </td>
				<td> <xsl:value-of select="MD_SIZE"/> </td>
			</tr>
		</xsl:for-each>
	</xsl:template>

	<xsl:template match="MD_SUMMARY">
		<p/>
		<table width="90%" rows="2" align="center">
			<tr align="center" bgcolor="skyblue">
				<th colspan="5"> Cumulative Summary for Memory Debugger Dumper Run </th>
			</tr>
			<tr>
				<td width="50%" bgcolor="ffffc0" align="center"> Valid Records </td>
				<td bgcolor="#e0f0d0"> <xsl:value-of select="MD_SUMMARY_VALID_RECORDS" /> </td>
			</tr>

			<tr>
				<td width="50%" bgcolor="ffffc0" align="center"> Invalid Records </td>
				<td bgcolor="#e0f0d0"> <xsl:value-of select="MD_SUMMARY_INVALID_RECORDS" /> </td>
			</tr>

			<tr>
				<td width="50%" bgcolor="ffffc0" align="center"> Total Number of Allocation/Deallocation Records </td>
				<td bgcolor="#e0f0d0"> <xsl:value-of select="MD_SUMMARY_TOTAL_RECORDS" /> </td>
			</tr>

		</table>
	</xsl:template>

	<xsl:template match="MD_FOOTER">
		<p/>
		<hr align="center" width="100%" color="maroon" />
		<h5 align="center"> <xsl:apply-templates/> </h5>
	</xsl:template>

</xsl:stylesheet>
