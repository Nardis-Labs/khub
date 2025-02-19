import React from 'react';
import { DataTable, IconButton, Table, TableBody, TableCell, TableContainer, TableHead, TableHeader, TableRow, TableToolbar, TableToolbarContent, TableToolbarSearch } from "@carbon/react";
import { AddFilled } from '@carbon/icons-react';

type AdminDataTableProps = {
  rows: any[];
  headers: any[];
  filterFunction: (args: any) => void;
  filterPlaceholder: string;
  filterValue: string;
  title: string;
  upsertFunction: (args: any) => void | undefined;
  upsertFunctionTitle: string | undefined;
};

export const AdminDataTable = (props: AdminDataTableProps) => {
  return (
    <DataTable rows={props.rows} headers={props.headers}>
      {({ rows, headers, getTableProps }) => (
        <TableContainer title={props.title}>
          <TableToolbar aria-label="admin data">
            <TableToolbarContent>
              <TableToolbarSearch value={props.filterValue} placeholder={props.filterPlaceholder} onChange={props.filterFunction} persistent />
              <IconButton kind="ghost" onClick={props.upsertFunction} label={props.upsertFunctionTitle}><AddFilled/></IconButton>
            </TableToolbarContent>
          </TableToolbar>
          <Table {...getTableProps()}>
            <TableHead>
              <TableRow>
                {headers.map((header) => (
                  <TableHeader key={header.key} className="theader">
                    {header.header}
                  </TableHeader>
                ))}
              </TableRow>
            </TableHead>
            <TableBody>
              {rows.map((row) => (
                <TableRow key={row.id}>
                  {row.cells.map((cell) => (
                    <TableCell key={cell.id}>{cell.value}</TableCell>
                  ))}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      )}
    </DataTable>
  );
};

