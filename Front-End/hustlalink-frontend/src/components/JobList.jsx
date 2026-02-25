export function JobCard({job}) {
  return (
    <div>
      <div>
        <p>{job.id}</p>
      </div>
       <div>
        <p>{job.name}</p>
      </div>
    </div>
  )
}


function JobList({ jobs }) {
  return (
    <div>
      {jobs.map(job => <JobCard key={job.id} job={job} />)}
    </div>
  );
}

export default JobList;