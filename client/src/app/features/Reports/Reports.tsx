import React from 'react';
import { ReportViewer } from '../../components/ReportViewer/ReportViewer';
import { useGetReportsQuery } from '../../../service/khub';
import { Loading } from '@carbon/react';

export const Reports = () => {

  const {data: reports = [], isLoading} = useGetReportsQuery({});

  return (
    <div style={{height: '100%'}}>
        {isLoading &&
          <div>
            <Loading withOverlay={true}/>
          </div>
        }
        <ReportViewer title='Pods' reports={reports}/>
    </div>
  );
};