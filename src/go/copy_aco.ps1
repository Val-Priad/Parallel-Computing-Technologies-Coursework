$aco_content = Get-Content aco.go
$two_opt = Get-Content two_opt.go

Set-Clipboard "
=== aco.go ===
$aco_content

=== two_opt ===
$two_opt
"
