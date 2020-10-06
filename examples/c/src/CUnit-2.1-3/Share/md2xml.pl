#!/usr/bin/perl
# This is in perl script for Memeory Debugging Report Generator

use warnings;
use strict;
use locale;

# Constant Declarations
use constant SIZE_INDEX => 0;
use constant POINTER_INDEX => 1;
use constant ALLOC_FILENAME_INDEX => 2;
use constant ALLOC_LINE_INDEX => 3;
use constant DEALLOC_FILENAME_INDEX => 4;
use constant DEALLOC_LINE_INDEX => 5;
use constant STATUS_INDEX => 6;

# Global Definitions
my @mem_list;

sub BuildDataList {
	
	open(LIST_FD, $ARGV[0]) or die "Failed to Open file $ARGV[0]";
	while (!eof(LIST_FD)) {
		my $line = <LIST_FD>;

		$line =~s/^[ \t]*//g;
		$line =~s/[ \t\n\r]*$//g;
		if (length($line) == 0 || index($line, "--") == 0) {
			next;
		}

		my @arr = split("\t", $line);
		my @alloc_arr = split(":", $arr[2]);
		my @dealloc_arr = split(":", $arr[3]);

		my $count = $#mem_list + 1;
		$mem_list[$count + 0] = $arr[0];
		$mem_list[$count + 1] = hex($arr[1]);
		$mem_list[$count + 2] = $alloc_arr[0];
		$mem_list[$count + 3] = $alloc_arr[1];
		$mem_list[$count + 4] = $dealloc_arr[0];
		$mem_list[$count + 5] = $dealloc_arr[1];
		$mem_list[$count + 6] = $arr[4];
		printf("%d  %d ", $mem_list[$count + SIZE_INDEX],
				$mem_list[$count + POINTER_INDEX])
	}

	close(LIST_FD);
	return 0;
}

sub PrintDataList {
	my $count = 0;
	while ($#mem_list > $count) {
#		printf("%d   ", $mem_list[$count + $SIZE_INDEX]);

		$count += STATUS_INDEX;

	}
}

BuildDataList();
PrintDataList();
