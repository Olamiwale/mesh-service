
echo "Testing Canary Deployment - 10% traffic to v2 (new version) and the 90% rto the v1 (old version)"
    
v1_count=0
v2_count=0
total_requests=100
    
for i in $(seq 1 $total_requests); do
response=$(curl -s http://product-service:8080/product/health)
version=$(echo $response | grep -o '"version":"v[12]"' | cut -d'"' -f4)
      
if [ "$version" == "v1" ]; then
    ((v1_count++))
    elif [ "$version" == "v2" ]; then
    ((v2_count++))
    fi
      
# Show progress
if [ $((i % 10)) -eq 0 ]; then
echo "Progress: $i/$total_requests requests"
fi
done
    
echo ""
echo "Results:"
echo "--------"
echo "v1 requests: $v1_count ($(( v1_count * 100 / total_requests ))%)"
echo "v2 requests: $v2_count ($(( v2_count * 100 / total_requests ))%)"
echo ""
echo "Expected: ~90% v1, ~10% v2"