// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-nocheck
import React from 'react';
import { CarbonIconType, Pagination, DataTable, Table, TableBatchActions, TableBatchAction, TableBody, TableCell, TableContainer, TableHead, TableHeader, TableRow, TableSelectRow, TableToolbar, TableToolbarContent, TableToolbarSearch, TableSelectAll } from "@carbon/react";

type ResourceDataTableProps = {
  rows: any[];
  headers: any[];
  filterFunction: (args: any) => void;
  filterPlaceholder: string;
  filterValue: string;
  title: string;
  batchActions: BatchActions[];
};

type BatchActions = {
  actionFunc: (args: any) => void;
  actionDescription: string;
  actionLabel: string;
  actionIcon: CarbonIconType;
}

export const ResourceDataTable = (props: ResourceDataTableProps) => {
 
  const [page, setPage] = React.useState(1);
  const [pageSize, setPageSize] = React.useState(20);

  const changePaginationState = (pageInfo) => {
    if (page != pageInfo.page) {
      setPage(pageInfo.page);
    }
    if (pageSize != pageInfo.pageSize) {
      setPageSize(pageInfo.pageSize);
    }
  };

  return (
    <>
    <DataTable rows={props.rows} headers={props.headers}>
      {({ rows, headers, getTableProps, getSelectionProps, getBatchActionProps, getHeaderProps, selectedRows}) => (
        
        <TableContainer title={props.title}>
          <TableToolbar aria-label="k8s resource data">
            <TableBatchActions {...getBatchActionProps()}>
              {props.batchActions.map((action) => (
                <TableBatchAction 
                  renderIcon={action.actionIcon} 
                  iconDescription={action.actionDescription} 
                  {...getBatchActionProps()} 
                  key={action.actionLabel} 
                  onClick={() => {
                    action.actionFunc(selectedRows);
                    setTimeout(() => { getBatchActionProps().onCancel(); }, 750);
                  }} 
                  tabIndex={getBatchActionProps().shouldShowBatchActions ? 0 : -1}>
                    {action.actionLabel}
                </TableBatchAction>
              ))}
            </TableBatchActions>
            <TableToolbarContent>
              <TableToolbarSearch value={props.filterValue} placeholder={props.filterPlaceholder} onChange={props.filterFunction} persistent />
            </TableToolbarContent>
          </TableToolbar>
          <Table {...getTableProps()}>
            <TableHead>
              <TableRow>
                <TableSelectAll {...getSelectionProps()} className="theader" disabled={selectedRows.length === 0}/>
                {headers.map((header) =>
                  <TableHeader key={header.key} className="theader" {...getHeaderProps({header, isSortable: header.isSortable})}>
                    {header.header}
                  </TableHeader>
                )}
              </TableRow>
            </TableHead>
            <TableBody>
              {
              rows.slice((page - 1) * pageSize).slice(0, pageSize).map((row) => (
                  <TableRow key={row.id}>
                    <TableSelectRow {...getSelectionProps({row})} disabled={props.batchActions.length === 0}/>
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
    <Pagination onChange={changePaginationState} page={page} pageSize={pageSize} pageSizes={[20, 100, 300]} totalItems={props.rows.length}/>
    </>
  );
};