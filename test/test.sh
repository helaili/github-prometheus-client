curl -H "Content-Type: application/json" -H "X-GitHub-Event: workflow_run" -d "@workflow/01_workflow_run_2096567160_requested.json" http://localhost:8080/webhook
curl -H "Content-Type: application/json" -H "X-GitHub-Event: workflow_job" -d "@workflow/02_workflow_job_5834689160_queued.json" http://localhost:8080/webhook
curl -H "Content-Type: application/json" -H "X-GitHub-Event: check_run" -d "@workflow/03_check_run_5834689160_created.json" http://localhost:8080/webhook
curl -H "Content-Type: application/json" -H "X-GitHub-Event: workflow_job" -d "@workflow/04_workflow_job_5834689160_completed.json" http://localhost:8080/webhook
curl -H "Content-Type: application/json" -H "X-GitHub-Event: check_run" -d "@workflow/05_check_run_5834689160_completed.json" http://localhost:8080/webhook
curl -H "Content-Type: application/json" -H "X-GitHub-Event: workflow_job" -d "@workflow/06_workflow_job_5834695996_queued.json" http://localhost:8080/webhook
curl -H "Content-Type: application/json" -H "X-GitHub-Event: check_run" -d "@workflow/07_check_run_5834695996_created.json" http://localhost:8080/webhook
curl -H "Content-Type: application/json" -H "X-GitHub-Event: deployment" -d "@workflow/08_deployment_540578780_created.json" http://localhost:8080/webhook
curl -H "Content-Type: application/json" -H "X-GitHub-Event: workflow_job" -d "@workflow/09_workflow_job_5834696229_queued.json" http://localhost:8080/webhook
curl -H "Content-Type: application/json" -H "X-GitHub-Event: check_run" -d "@workflow/10_check_run_5834696229_created.json" http://localhost:8080/webhook
curl -H "Content-Type: application/json" -H "X-GitHub-Event: workflow_job" -d "@workflow/11_workflow_job_5834695996_in_progress.json" http://localhost:8080/webhook
curl -H "Content-Type: application/json" -H "X-GitHub-Event: workflow_job" -d "@workflow/12_workflow_job_5834696229_in_progress.json" http://localhost:8080/webhook
curl -H "Content-Type: application/json" -H "X-GitHub-Event: deployment_status" -d "@workflow/13_deployment_status_540578780_created.json" http://localhost:8080/webhook
curl -H "Content-Type: application/json" -H "X-GitHub-Event: check_run" -d "@workflow/14_check_run_5834695996_completed.json" http://localhost:8080/webhook
curl -H "Content-Type: application/json" -H "X-GitHub-Event: workflow_job" -d "@workflow/15_workflow_job_5834695996_completed.json" http://localhost:8080/webhook
curl -H "Content-Type: application/json" -H "X-GitHub-Event: workflow_job" -d "@workflow/16_workflow_job_5834696229_completed.json" http://localhost:8080/webhook
curl -H "Content-Type: application/json" -H "X-GitHub-Event: check_run" -d "@workflow/17_check_run_5834696229_completed.json" http://localhost:8080/webhook
curl -H "Content-Type: application/json" -H "X-GitHub-Event: deployment_status" -d "@workflow/18_deployment_status_540578780_created.json" http://localhost:8080/webhook
curl -H "Content-Type: application/json" -H "X-GitHub-Event: check_suite" -d "@workflow/19_check_suite_5939269259_completed.json" http://localhost:8080/webhook
curl -H "Content-Type: application/json" -H "X-GitHub-Event: workflow_run" -d "@workflow/20_workflow_run_2096567160_completed.json" http://localhost:8080/webhook
curl -H "Content-Type: application/json" -H "X-GitHub-Event: workflow_run" -d "@workflow/101_workflow_run_2096567160_requested.json" http://localhost:8080/webhook
