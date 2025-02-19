import { Button } from '@carbon/react';
import React from 'react';
import { ResourceDataTable } from '../ResourceDataTable/ResourceDataTable';
import { FaDownload } from "react-icons/fa6";
import { khubApi } from '../../../service/khub';
import { useAppDispatch } from '../../store';

type DocViewerProps = {
  title: string;
  reports: any[];
};


export const ReportViewer: React.FC<DocViewerProps> = (props) => {

  const dispatch = useAppDispatch();
  const [reportFilter, setReportFilter] = React.useState('');
  const filterReports = (args: any) => {
    setReportFilter(args.target.value);
  };

  const handleDownloadReport = (report: string) => {
    const result = dispatch(khubApi.endpoints.getReportDownloadURL.initiate({key: report}));
   
    result.then((response) => {
      if (response !== undefined && response.data !== "") {
      const win = window.open(response.data);
      if (win !== null && !report.includes('.log')) {
        win.focus();
      }
    }
    });
  };

  const headers = [
    {
      header: 'Name',
      key: 'name',
      isSortable: true
    },
    {
      header: 'Type',
      key: 'type',
      isSortable: true
    },
    {
      header: 'Created By',
      key: 'createdBy',
      isSortable: true
    },
    {
      header: 'Reference',
      key: 'reference',
      isSortable: true
    },
    {
      header: 'Created Date',
      key: 'created'
    },
    {
      header: 'Size',
      key: 'size',
      isSortable: true
    },
    {
      header: '',
      key: 'controls'
    }
  ];

  const rows: any[] = props.reports.filter((report: any) => {
    return (
      (report?.name !== "") && // Sometimes these files dont have names, so we dont want to display them.
      ((report?.name !== undefined && report.name.toLowerCase().includes(reportFilter.toLowerCase())) ||
      (report?.type !== undefined && report.type.toLowerCase().includes(reportFilter.toLowerCase())) ||
      (report?.user !== undefined && report.user.toLowerCase().includes(reportFilter.toLowerCase())) ||
      (report?.reason !== undefined && report.reason.toLowerCase().includes(reportFilter.toLowerCase())))
    );
  }).map((report: any) => {
    return {
      id: report.name,
      name: report.name,
      type: report.type,
      createdBy: report.user,
      reference: report.reason,
      created: report.created,
      size: report.size + report.sizeUnits,
      controls: <div style={{float: 'right'}}>
                   <Button onClick={() => handleDownloadReport(report.name)} renderIcon={FaDownload} kind='ghost' hasIconOnly iconDescription='Download'/>
                </div>
    };
  });

  return (
    <div id={'ReportViewer'} className="ReportViewer">
      <ResourceDataTable 
        rows={rows} 
        headers={headers} 
        filterFunction={(filterReports)} 
        filterPlaceholder={'Filter reports'}
        filterValue={reportFilter}
        title={'Reports (' + props.reports.length + ')'}
        batchActions={[]}
      />
    </div>
  );
};
