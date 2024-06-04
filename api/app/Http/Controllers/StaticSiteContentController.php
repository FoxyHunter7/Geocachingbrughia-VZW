<?php

namespace App\Http\Controllers;

use App\Services\StaticSiteContentService;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

class StaticSiteContentController
{
    private $_service;

    public function __construct(StaticSiteContentService $service)
    {
        $this->_service = $service;
    }

    public function all(Request $request, $isAdmin)
    {
        $search = $request->query('search', '');
        $sortBy = $request->query('sort_by', '');
        $sortDirection = $request->query('sort_direction', 'desc');

        return $this->_service->all(search: $search, sortBy: $sortBy, sortDirection: $sortDirection, isAdmin: $isAdmin);
    }

    public function allPublic (Request $request)
    {
        return response()->json($this->all($request, false), Response::HTTP_OK);
    }

    public function allAdmin (Request $request)
    {
        return response()->json($this->all($request, true), Response::HTTP_OK);
    }

    public function add(Request $request)
    {
        $data = $request->all();
        $staticSiteContent = $this->_service->add($data);

        if ($this->_service->hasErrors()) {
            return [
                "errors" => $this->_service->getErrors()
            ];
        }

        return ["data" => $staticSiteContent];
    }

    public function update(Request $request, string $property)
    {
        $data = $request->all();
        $staticSiteContent = $this->_service->update($property, $data);

        if ($this->_service->hasErrors()) {
            return [
                "errors" => $this->_service->getErrors()
            ];
        }

        return ["data" => $staticSiteContent];
    }

    public function delete(string $property)
    {
        $deletedStaticSiteContent = $this->_service->delete($property);

        if ($this->_service->hasErrors()) {
            return [
                "errors" => $this->_service->getErrors()
            ];
        }

        return ["deleted" => $deletedStaticSiteContent];
    }
}
