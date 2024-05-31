<?php

namespace Database\Seeders;

use App\Models\StaticSiteProperty;
use Illuminate\Database\Seeder;

class StaticSitePropertiesSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        $data = ReadFromCsv::getDataFromCsv(storage_path('app/data/csv/staticSiteProperties.csv'), ';');

        $model = new StaticSiteProperty();
        foreach ($data as $row) {
            $model->create($row);
        }
    }
}
