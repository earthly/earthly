       IDENTIFICATION DIVISION.
           PROGRAM-ID. "Fibonacci".
       ENVIRONMENT DIVISION.
       DATA DIVISION.
       WORKING-STORAGE SECTION.
       01  ix                    BINARY-C-LONG VALUE 0.
       01  first-number          BINARY-C-LONG VALUE 0.
       01  second-number         BINARY-C-LONG VALUE 1.
       01  temp-number           BINARY-C-LONG VALUE 1.
       01  display-number        PIC Z(3)9.
       PROCEDURE DIVISION.
       START-PROGRAM.
           MOVE first-number TO display-number.
           DISPLAY display-number.
           MOVE second-number TO display-number.
           DISPLAY display-number.
           PERFORM VARYING ix FROM 1 BY 1 UNTIL ix = 10
               ADD first-number TO second-number GIVING temp-number
               MOVE second-number TO first-number
               MOVE temp-number TO second-number
               MOVE temp-number TO display-number
               DISPLAY display-number
           END-PERFORM.
           STOP RUN.
